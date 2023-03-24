package mongo

import (
	"context"
	"fmt"
	config "github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

type CreateRequestPayload struct {
	CollectionName string
	Payload        interface{}
}

type RetrieveRequestPayload struct {
	CollectionName string
	Id             string
	Filters        bson.M
}

func Connection(ctx context.Context) (*Mongo, errors.Error) {
	lock.Lock()
	defer lock.Unlock()
	if Client == nil {
		mongoClient, err := _initConnection(ctx)
		if !err.Equals(errors.NotFound) {
			return nil, err
		}
		Client = mongoClient
	}
	return Client, errors.OK
}

func _initConnection(ctx context.Context) (*Mongo, errors.Error) {
	config := config.Current.Mongo
	client, err := mongo.NewClient(options.Client().ApplyURI(
		"mongodb://" + config.Username + ":" + config.Password + "@" + config.Address + ":" + strconv.Itoa(config.Port)),
	)
	if err != nil {
		logger.Error.Printf("Error creating mongo client: %s", err)
		return &Mongo{}, errors.ExternalServiceError
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabases(ctx, bson.M{})
	if err != nil {
		msg := fmt.Sprintf("Error listing databases: %s", err)
		logger.Error.Println(msg)
		return &Mongo{Client: client}, errors.ExternalServiceError
	}
	for _, db := range databases.Databases {
		if db.Name == config.Database {
			matchDB := client.Database(config.Database)

			return &Mongo{Client: client, Database: matchDB},
				errors.ExternalServiceError.WithMessage(
					option.New(reflect.TypeOf(err).String(), err),
				)
		}
	}
	return &Mongo{client, &mongo.Database{}}, errors.OK
}

// Create args - context.Context - Payload interface{}
func (m *Mongo) Create(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	if !optn.SetType(reflect.TypeOf(CreateRequestPayload{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(CreateRequestPayload{}).Kind(), optn.Value,
			),
		)
	}
	payload := optn.Value.(CreateRequestPayload)
	collection := m.Database.Collection(payload.CollectionName)
	result, err := collection.InsertOne(ctx, payload.Payload)
	if err != nil {
		return nil, errors.ExternalServiceError.WithMessage(err)
	}
	return result.InsertedID, errors.OK
}

func (m *Mongo) Get(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	// Collection string - Payload interface{} - Filter interface{}
	if !optn.SetType(reflect.TypeOf(RetrieveRequestPayload{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(RetrieveRequestPayload{}).Kind(), optn.Value,
			),
		)
	}
	var res interface{}
	payload := optn.Value.(RetrieveRequestPayload)
	searchWithFilter := payload.Filters != nil
	// Collection string - Payload interface{}
	searchById := payload.Id != ""
	collection := m.Database.Collection(payload.CollectionName)

	if searchWithFilter {
		cursor, err := collection.Find(ctx, payload.Filters)
		if err != nil {
			return nil, errors.NotFound.WithMessage(err)
		}
		if err = cursor.All(ctx, &res); err != nil {
			return nil, errors.ExternalServiceError.WithMessage(err)
		}
	} else if searchById {
		cursor := collection.FindOne(ctx, bson.M{"_id": payload.Id})
		if err := cursor.Decode(&res); err != nil {
			return nil, errors.NotFound.WithMessage(err)
		}
	} else {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(RetrieveRequestPayload{}).Kind(), optn.Value,
			),
		)
	}
	return res, errors.OK
}

func (m *Mongo) GetAll(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	if !optn.SetType(reflect.TypeOf(RetrieveRequestPayload{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(RetrieveRequestPayload{}).Kind(), optn.Value,
			),
		)
	}
	var res []interface{}
	payload := optn.Value.(RetrieveRequestPayload)
	searchWithFilter := payload.Filters == nil
	collection := m.Database.Collection(payload.CollectionName)

	filter := bson.M{}
	if searchWithFilter {
		filter = payload.Filters
	}
	curs, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.NotFound.WithMessage(err)
	}
	if err = curs.All(ctx, &res); err != nil {
		return nil, errors.ExternalServiceError.WithMessage(err)
	}

	return nil,
		errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid arguments: must be (string, []interface{}, bson.M) or (string, []interface{}), and not %v",
				optn.Value,
			),
		)
}

func (m *Mongo) Update(ctx context.Context, optn option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (m *Mongo) Delete(ctx context.Context, optn option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (m *Mongo) Close(ctx context.Context) errors.Error {
	err := m.Client.Disconnect(ctx)
	if err != nil {
		return errors.ExternalServiceError.WithMessage(err)
	}
	return errors.OK
}
