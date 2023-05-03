package mongo

import (
	"context"
	"fmt"
	config "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strconv"
	"sync"
)

var Client *DAO
var lock = &sync.Mutex{}

type DAO struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	if Client == nil {
		logger.Info.Printf("Init mongo client")
		conf := config.Current.Mongo
		client, err := mongo.NewClient(options.Client().ApplyURI(
			"mongodb://" + conf.Username + ":" + conf.Password + "@" + conf.Address + ":" + strconv.Itoa(conf.Port)),
		)
		if err != nil {
			logger.Error.Printf("Error creating mongo client: %s", err)
			panic(err)
		}
		err = client.Connect(context.Background())
		if err != nil {
			logger.Error.Printf("Error connecting to mongo: %s", err)
			panic(err)
		}
		matchDB := client.Database(conf.Database)
		Client = &DAO{
			Client:   client,
			Database: matchDB,
		}
	}
}

// Create args - context.Context - Payload interface{}
func (m *DAO) Create(ctx context.Context, opt option.Option) interface{} {
	if !opt.SetType(reflect.TypeOf(CreateRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(CreateRequest{}).Kind(), opt.Value,
			),
		))
	}
	payload := opt.Value.(CreateRequest)
	collection := m.Database.Collection(payload.CollectionName)
	result, err := collection.InsertOne(ctx, payload.Payload)
	if err != nil {
		panic(errors.ExternalServiceError.WithMessage(err))
	}
	objectID := result.InsertedID.(primitive.ObjectID)

	return CreateResponse{Id: objectID.Hex()}
}

func (m *DAO) Get(ctx context.Context, opt option.Option) interface{} {
	// Collection string - Payload interface{} - Filter interface{}
	if !opt.SetType(reflect.TypeOf(GetRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(GetRequest{}).Kind(), opt.Value,
			),
		))
	}
	payload := opt.Value.(GetRequest)
	collection := m.Database.Collection(payload.CollectionName)
	id, err := primitive.ObjectIDFromHex(payload.Id)
	if err != nil {
		panic(errors.InvalidArgument.WithMessage(err.Error()))
	}
	filter := bson.M{"_id": id}
	res, err := collection.FindOne(ctx, filter).DecodeBytes()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			panic(errors.NotFound.WithMessage(fmt.Sprintf("No documents found for id %s", payload.Id)))
		} else {
			panic(errors.ExternalServiceError.WithMessage(err.Error()))
		}
	}
	return GetResponse{Payload: res}
}

func (m *DAO) Exists(ctx context.Context, opt option.Option) bool {
	if !opt.SetType(reflect.TypeOf(ExistsRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(ExistsRequest{}).Kind(), opt.Value,
			),
		))
	}
	payload := opt.Value.(ExistsRequest)
	collection := m.Database.Collection(payload.CollectionName)
	count, err := collection.CountDocuments(ctx, payload.Filter)
	if err != nil {
		panic(errors.ExternalServiceError.WithMessage(err.Error()))
	}
	return count > 0
}

func (m *DAO) GetAll(ctx context.Context, optn option.Option) interface{} {
	if !optn.SetType(reflect.TypeOf(GetAllRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(GetRequest{}).Kind(), optn.Value,
			),
		))
	}
	payload := optn.Value.(GetAllRequest)
	collection := m.Database.Collection(payload.CollectionName)

	filter := payload.Filter
	curs, err := collection.Find(ctx, filter)
	if err != nil {
		panic(errors.NotFound.WithMessage(err))
	}
	defer func(curs *mongo.Cursor, ctx context.Context) {
		errCurs := curs.Close(ctx)
		if errCurs != nil {
			logger.Error.Printf("Error closing cursor: %s", errCurs)
		}
	}(curs, ctx)
	var res []interface{}
	for curs.Next(ctx) {
		var elem interface{}
		errCur := curs.Decode(&elem)
		if errCur != nil {
			panic(errors.ExternalServiceError.WithMessage(errCur))
		}
		res = append(res, elem)
	}
	return GetAllResponse{Payload: res}
}

func (m *DAO) Update(ctx context.Context, opt option.Option) {
	if !opt.SetType(reflect.TypeOf(UpdateRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(UpdateRequest{}).Kind(), opt.Value,
			),
		))
	}
	payload := opt.Value.(UpdateRequest)
	collection := m.Database.Collection(payload.CollectionName)
	bsonMap := parseBsonMap(structs.Map(payload.Payload))
	if _, err := collection.UpdateByID(ctx, payload.Id, bson.M{"$set": bsonMap}); err != nil {
		panic(errors.ExternalServiceError.WithMessage(err.Error()))
	}
}

func (m *DAO) Delete(ctx context.Context, opt option.Option) {
	if !opt.SetType(reflect.TypeOf(DeleteRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(DeleteRequest{}).Kind(), opt.Value,
			),
		))
	}
	payload := opt.Value.(DeleteRequest)
	collection := m.Database.Collection(payload.CollectionName)
	id, err := primitive.ObjectIDFromHex(payload.Id)
	if err != nil {
		panic(errors.InvalidArgument.WithMessage(err.Error()))
	}
	res := collection.FindOneAndDelete(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		panic(errors.NotFound.WithMessage(res.Err().Error()))
	}
}

func (m *DAO) Close(ctx context.Context) {
	err := m.Client.Disconnect(ctx)
	if err != nil {
		panic(errors.ExternalServiceError.WithMessage(err))
	}
}
