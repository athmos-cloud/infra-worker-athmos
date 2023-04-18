package mongo

import "go.mongodb.org/mongo-driver/bson"

type CreateRequest struct {
	CollectionName string
	Payload        interface{}
}

type CreateResponse struct {
	Id string
}

type GetRequest struct {
	CollectionName string
	Id             string
	Payload        interface{}
}

type ExistsRequest struct {
	CollectionName string
	Filter         bson.M
}

type GetResponse struct {
	Payload bson.Raw
}

type GetAllRequest struct {
	CollectionName string
	Filter         bson.M
}

type GetAllResponse struct {
	Payload []interface{}
}

type UpdateRequest struct {
	CollectionName string
	Id             string
	Payload        interface{}
}

type DeleteRequest struct {
	CollectionName string
	Id             string
}