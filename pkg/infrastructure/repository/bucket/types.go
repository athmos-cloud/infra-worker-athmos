package bucket

import (
	"github.com/minio/minio-go/v7"
)

const (
	_DefaultMinioRegion = ""
	_DefaultUseSSL      = false
	_DefaultObjectName  = "/"
	_RemoveBucketFlag   = "@remove"
)

type Minio struct {
	Client *minio.Client
}

type CreateRequestPayload struct {
	BucketName string
	FromDir    string
	ToDir      string
}

type UpdateRequestPayload struct {
	BucketName string
	FromDir    string
	ToDir      string
}

type RetrieveRequestPayload struct {
	BucketName string
	Dir        string
}

type DeleteRequestPayload struct {
	BucketName string
	Dir        string
}

type Bucket struct {
	Name   string
	Client *Minio
}
