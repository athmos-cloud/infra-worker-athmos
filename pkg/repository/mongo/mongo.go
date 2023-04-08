package mongo

import (
	"context"
	"fmt"
	config "github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	internalCtx "github.com/PaulBarrie/infra-worker/pkg/kernel/context"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strconv"
	"sync"
)

var Client *Mongo
var lock = &sync.Mutex{}

type Mongo struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	if Client == nil {
		config := config.Current.Mongo
		client, err := mongo.NewClient(options.Client().ApplyURI(
			"mongodb://" + config.Username + ":" + config.Password + "@" + config.Address + ":" + strconv.Itoa(config.Port)),
		)
		if err != nil {
			logger.Error.Printf("Error creating mongo client: %s", err)
			panic(err)
		}
		err = client.Connect(internalCtx.Current)
		if err != nil {
			logger.Error.Printf("Error connecting to mongo: %s", err)
			panic(err)
		}
		matchDB := client.Database(config.Database)
		Client = &Mongo{
			Client:   client,
			Database: matchDB,
		}
	}
}

// Create args - context.Context - Payload interface{}
func (m *Mongo) Create(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	if !optn.SetType(reflect.TypeOf(CreateRequest{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(CreateRequest{}).Kind(), optn.Value,
			),
		)
	}
	payload := optn.Value.(CreateRequest)
	collection := m.Database.Collection(payload.CollectionName)
	result, err := collection.InsertOne(ctx, payload.Payload)
	if err != nil {
		return nil, errors.ExternalServiceError.WithMessage(err)
	}
	objectID := result.InsertedID.(primitive.ObjectID)

	return CreateResponse{Id: objectID.Hex()}, errors.OK
}

func (m *Mongo) Get(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	// Collection string - Payload interface{} - Filter interface{}
	if !optn.SetType(reflect.TypeOf(GetRequest{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(GetRequest{}).Kind(), optn.Value,
			),
		)
	}
	payload := optn.Value.(GetRequest)
	collection := m.Database.Collection(payload.CollectionName)
	id, err := primitive.ObjectIDFromHex(payload.Id)
	if err != nil {
		return nil, errors.InvalidArgument.WithMessage(err.Error())
	}
	filter := bson.M{"_id": id}
	res, err := collection.FindOne(ctx, filter).DecodeBytes()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info.Printf("No documents found for id %s", payload.Id)
			return nil, errors.NotFound
		} else {
			return nil, errors.ExternalServiceError.WithMessage(err.Error())
		}
	}
	return GetResponse{Payload: res}, errors.OK
}

func (m *Mongo) GetAll(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	if !optn.SetType(reflect.TypeOf(GetAllRequest{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(GetRequest{}).Kind(), optn.Value,
			),
		)
	}
	payload := optn.Value.(GetAllRequest)
	collection := m.Database.Collection(payload.CollectionName)

	filter := payload.Filter
	curs, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.NotFound.WithMessage(err)
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
			return nil, errors.ExternalServiceError.WithMessage(errCur)
		}
		res = append(res, elem)
	}
	return GetAllResponse{Payload: res}, errors.OK
}

func (m *Mongo) Update(ctx context.Context, optn option.Option) errors.Error {
	if !optn.SetType(reflect.TypeOf(UpdateRequest{}).String()).Validate() {
		return errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(UpdateRequest{}).Kind(), optn.Value,
			),
		)
	}
	payload := optn.Value.(UpdateRequest)
	collection := m.Database.Collection(payload.CollectionName)
	if _, err := collection.UpdateByID(ctx, payload.Id, payload.Payload); err != nil {
		return errors.ExternalServiceError.WithMessage(err)
	}
	return errors.OK
}

func (m *Mongo) Delete(ctx context.Context, optn option.Option) errors.Error {
	if !optn.SetType(reflect.TypeOf(DeleteRequest{}).String()).Validate() {
		return errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(DeleteRequest{}).Kind(), optn.Value,
			),
		)
	}
	payload := optn.Value.(DeleteRequest)
	collection := m.Database.Collection(payload.CollectionName)
	res := collection.FindOneAndDelete(ctx, bson.M{"_id": payload.Id})
	if res.Err() != nil {
		return errors.NotFound.WithMessage(res.Err().Error())
	}
	return errors.OK
}

func (m *Mongo) Close(ctx context.Context) errors.Error {
	err := m.Client.Disconnect(ctx)
	if err != nil {
		return errors.ExternalServiceError.WithMessage(err)
	}
	return errors.OK
}
