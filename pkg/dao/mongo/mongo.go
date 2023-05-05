package mongo

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"sync"
)

var Client *DAO
var lock = &sync.Mutex{}

type DAO struct{}

var inited = false

func init() {
	lock.Lock()
	defer lock.Unlock()
	if !inited {
		logger.Info.Printf("Init mongo client")
		conf := config.Current.Mongo
		uri := "mongodb://" + conf.Username + ":" + conf.Password + "@" + conf.Address + ":" + strconv.Itoa(conf.Port)
		if err := mgm.SetDefaultConfig(nil, conf.Database, options.Client().ApplyURI(uri)); err != nil {
			panic(errors.InternalError.WithMessage(err.Error()))
		}
		inited = true
	}
}

//// Create args - context.Context - Payload interface{}
//func (m *DAO) Create(ctx context.Context, opt option.Option) interface{} {
//	if !opt.SetType(reflect.TypeOf(CreateRequest{}).String()).Validate() {
//		panic(errors.InvalidArgument.WithMessage(
//			fmt.Sprintf(
//				"Invalid argument type, expected %s, got %v", reflect.TypeOf(CreateRequest{}).Kind(), opt.Value,
//			),
//		))
//	}
//	request := opt.Value.(CreateRequest)
//	collection := mgm.Coll(request.Payload).Collection
//
//	return CreateResponse{Id: objectID.Hex()}
//}
//
//func (m *DAO) Get(ctx context.Context, opt option.Option) interface{} {
//	// Collection string - Payload interface{} - Filter interface{}
//	if !opt.SetType(reflect.TypeOf(GetRequest{}).String()).Validate() {
//		panic(errors.InvalidArgument.WithMessage(
//			fmt.Sprintf(
//				"Invalid argument type, expected %s, got %v", reflect.TypeOf(GetRequest{}).Kind(), opt.Value,
//			),
//		))
//	}
//	payload := opt.Value.(GetRequest)
//
//	return GetResponse{Payload: res}
//}
//
//func (m *DAO) Exists(ctx context.Context, opt option.Option) bool {
//	if !opt.SetType(reflect.TypeOf(ExistsRequest{}).String()).Validate() {
//		panic(errors.InvalidArgument.WithMessage(
//			fmt.Sprintf(
//				"Invalid argument type, expected %s, got %v", reflect.TypeOf(ExistsRequest{}).Kind(), opt.Value,
//			),
//		))
//	}
//	payload := opt.Value.(ExistsRequest)
//
//	return count > 0
//}
//
//func (m *DAO) GetAll(ctx context.Context, optn option.Option) interface{} {
//	if !optn.SetType(reflect.TypeOf(GetAllRequest{}).String()).Validate() {
//		panic(errors.InvalidArgument.WithMessage(
//			fmt.Sprintf(
//				"Invalid argument type, expected %s, got %v", reflect.TypeOf(GetRequest{}).Kind(), optn.Value,
//			),
//		))
//	}
//	payload := optn.Value.(GetAllRequest)
//
//	return GetAllResponse{Payload: res}
//}
//
//func (m *DAO) Update(ctx context.Context, opt option.Option) {
//	if !opt.SetType(reflect.TypeOf(UpdateRequest{}).String()).Validate() {
//		panic(errors.InvalidArgument.WithMessage(
//			fmt.Sprintf(
//				"Invalid argument type, expected %s, got %v", reflect.TypeOf(UpdateRequest{}).Kind(), opt.Value,
//			),
//		))
//	}
//	request := opt.Value.(UpdateRequest)
//
//}
//
//func (m *DAO) Delete(ctx context.Context, opt option.Option) {
//	if !opt.SetType(reflect.TypeOf(DeleteRequest{}).String()).Validate() {
//		panic(errors.InvalidArgument.WithMessage(
//			fmt.Sprintf(
//				"Invalid argument type, expected %s, got %v", reflect.TypeOf(DeleteRequest{}).Kind(), opt.Value,
//			),
//		))
//	}
//	payload := opt.Value.(DeleteRequest)
//
//}
//
//func (m *DAO) Close(ctx context.Context) {
//}
