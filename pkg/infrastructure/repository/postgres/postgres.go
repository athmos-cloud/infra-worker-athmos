package postgres

import (
	"context"
	"fmt"
	config2 "github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"sync"
)

var Client *Repository
var lock = &sync.Mutex{}

const DefaultSQLDialect = "postgres"

type Repository struct {
	Database *gorm.DB
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

func Connection(ctx context.Context) (*Repository, errors.Error) {
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

func _initConnection(ctx context.Context) (*Repository, errors.Error) {
	config := config2.Get().Postgres
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Address, config.Port, config.Username, config.Password, config.Database, config.SSLMode,
	)
	db, err := gorm.Open(DefaultSQLDialect, psqlInfo)
	if err != nil {
		panic("failed to connect database")
	}

	return &Repository{db}, errors.OK
}

// Create args - context.Context - Payload interface{}
func (m *Repository) Create(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	if !optn.SetType(reflect.TypeOf(CreateRequestPayload{}).Kind()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(CreateRequestPayload{}).Kind(), optn.Value,
			),
		)
	}
	return nil, errors.OK
}

func (m *Repository) Get(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	// Collection string - Payload interface{} - Filter interface{}
	if !optn.SetType(reflect.TypeOf(RetrieveRequestPayload{}).Kind()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(RetrieveRequestPayload{}).Kind(), optn.Value,
			),
		)
	}

	return nil, errors.OK
}

func (m *Repository) GetAll(ctx context.Context, optn option.Option) ([]interface{}, errors.Error) {
	if !optn.SetType(reflect.TypeOf(RetrieveRequestPayload{}).Kind()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(RetrieveRequestPayload{}).Kind(), optn.Value,
			),
		)
	}

	return nil,
		errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid arguments: must be (string, []interface{}, bson.M) or (string, []interface{}), and not %v",
				optn.Value,
			),
		)
}

func (m *Repository) Update(ctx context.Context, optn option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (m *Repository) Delete(ctx context.Context, optn option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (m *Repository) Close(_ context.Context) errors.Error {
	err := m.Database.Close()
	if err != nil {
		return errors.ExternalServiceError.WithMessage(err)
	}
	return errors.OK
}
