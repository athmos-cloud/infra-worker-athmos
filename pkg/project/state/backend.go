package state

type BackendType string

const (
	BackendTypeMinio BackendType = "minio"
)

type Backend struct {
	Repository interface{} `bson:"repository"`
}

type IBackend interface {
	Get() interface{}
	Generate() error
}
