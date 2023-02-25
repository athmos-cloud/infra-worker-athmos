package bucket

import (
	"context"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/types/workdir"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

var MinioClient = &Minio{}
var lock = &sync.Mutex{}

func Connection() (*Minio, errors.Error) {
	if MinioClient == nil || reflect.ValueOf(MinioClient).IsNil() {

		minioConnection, err := _initConnection()
		if !err.IsOk() {
			return nil, err
		}
		lock.Lock()
		defer lock.Unlock()
		MinioClient = minioConnection
	}
	return MinioClient, errors.OK
}

func _initConnection() (*Minio, errors.Error) {
	var err error
	conf := config.Get().Minio

	minioClient, err := minio.New(config.Get().Minio.Address, &minio.Options{
		Creds: credentials.NewStaticV4(
			conf.AccessKeyID,
			conf.SecretAccessKey,
			conf.Token,
		),
		Secure: conf.UseSSL,
		Region: conf.Region,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return &Minio{
		Client: minioClient,
	}, errors.OK
	// Initialize minio client object.

}

func (m *Minio) Create(context context.Context, optn option.Option) (interface{}, errors.Error) {
	if optn.SetType(reflect.String); !optn.Validate() {
		return nil, errors.InvalidArgument.WithMessage("options are not valid : must be a string bucket name")
	}
	bucketName := optn.Value.(string)
	err := MinioClient.Client.MakeBucket(context, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := MinioClient.Client.BucketExists(context, bucketName)
		if errBucketExists == nil && exists {
			logger.Info.Println("bucket %s already ownd \n", bucketName)
		} else {
			logger.Error.Println(err)
			return nil,
				errors.ExternalServiceError.WithMessage(
					fmt.Sprintf("error creating bucket %s : %v", bucketName, err),
				)
		}
	} else {
		logger.Info.Println("Successfully created bucket %s\n", bucketName)
	}
	return nil, errors.OK
}

// Get 's args must be a name of the bucket, name of the object to get in it
func (m *Minio) Get(context context.Context, optn option.Option) (interface{}, errors.Error) {
	if optn.SetType(reflect.TypeOf(RetrieveRequestPayload{}).Kind()); !optn.Validate() {
		return nil,
			errors.InvalidArgument.WithMessage(
				fmt.Sprintf(
					"options are not valid : must be a string bucket name and a string object name, got %v", optn.Value,
				),
			)
	}
	payload := optn.Value.(RetrieveRequestPayload)

	isFile := func(objectName string) bool {
		return objectName[len(objectName)-1:] != "/"
	}

	objToString := func(object *minio.Object) (string, error) {
		content, err := object.Stat()
		if err != nil {
			logger.Error.Println(err)
			return "", err
		}
		contentBytes := make([]byte, content.Size)
		_, err = object.Read(contentBytes)
		if err != nil {
			logger.Error.Println(err)
			return "", err
		}
		return string(contentBytes[:]), nil
	}

	var resFolder workdir.Folder
	if isFile(payload.Dir) {
		// Get the object
		object, err := m.Client.GetObject(context, payload.BucketName, payload.Dir, minio.GetObjectOptions{})
		if err != nil {
			logger.Error.Println(err)
			return nil, errors.ExternalServiceError.WithMessage(err.Error())
		}
		res, err := objToString(object)
		if err != nil {
			logger.Error.Println(err)
			return nil, errors.ConversionError.WithMessage(err.Error())
		}
		resFolder.Files = append(resFolder.Files, workdir.File{Name: payload.Dir, Content: res})
		return resFolder, errors.OK
	} else {
		// Get the object
		objects := m.Client.ListObjects(context, payload.BucketName, minio.ListObjectsOptions{Prefix: payload.Dir})
		var files []workdir.File

		for object := range objects {
			if object.Err != nil {
				logger.Error.Println(object.Err)
				return resFolder, errors.ExternalServiceError.WithMessage(object.Err.Error())
			}
			if isFile(object.Key) {
				fileFolder, err := m.Get(
					context,
					option.Option{
						Type: reflect.TypeOf(RetrieveRequestPayload{}).Kind(),
						Value: RetrieveRequestPayload{
							BucketName: payload.BucketName,
							Dir:        object.Key,
						},
					},
				)
				if !err.IsOk() {
					logger.Error.Println(err)
					return resFolder, err
				}
				files = append(files, workdir.File{Name: object.Key, Content: fileFolder.(workdir.Folder).Files[0].Content})
			} else {
				folderObject, err := m.Get(
					context,
					option.Option{
						Type: reflect.TypeOf(RetrieveRequestPayload{}).Kind(),
						Value: RetrieveRequestPayload{
							BucketName: payload.BucketName,
							Dir:        object.Key,
						}},
				)
				if !err.IsOk() {
					logger.Error.Println(err)
				}
				resFolder.Folders = append(resFolder.Folders, folderObject.(workdir.Folder))
			}
		}
		return resFolder, errors.OK
	}
}

func (m *Minio) GetAll(ctx context.Context, bucketName option.Option) ([]interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}

// Update Args must be a name of the bucket, name of the object to get in it
func (m *Minio) Update(ctx context.Context, optn option.Option) errors.Error {
	if !optn.SetType(reflect.TypeOf(UpdateRequestPayload{}).Kind()).Validate() {
		return errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"options are not valid : must be a string bucket name and a string object name, got %v", optn,
			),
		)
	}
	payload := optn.Value.(UpdateRequestPayload)

	err := filepath.Walk(payload.FromDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		// Create a new object in the MinIO bucket for each file in the local folder
		objectName := filepath.Join(payload.ToDir, path[len(payload.FromDir)+1:])
		contentType := "application/octet-stream"
		reader, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = m.Client.PutObject(context.Background(), payload.BucketName, payload.ToDir, reader, info.Size(), minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			return err
		}

		logger.Info.Printf("Uploaded %s to %s\n", path, objectName)
		return nil
	})
	if err != nil {
		logger.Error.Printf("Error while uploading file %s to bucket %s: %v", payload.FromDir, payload.BucketName, err)
		return errors.ExternalServiceError.WithMessage(err.Error())
	}
	return errors.OK
}

func (m *Minio) Delete(ctx context.Context, option option.Option) errors.Error {
	if !option.SetType(reflect.TypeOf(DeleteRequestPayload{}).Kind()).Validate() {
		return errors.InvalidArgument.WithMessage("args must be a string bucket name and an optional string object name")
	}
	payload := option.Value.(DeleteRequestPayload)
	removeBucket := !(payload.Dir != _RemoveBucketFlag)
	if removeBucket {
		err := m.Client.RemoveBucket(ctx, payload.BucketName)
		if err != nil {
			logger.Error.Printf("Error while deleting bucket %args: %args", payload.BucketName, err)
			return errors.ExternalServiceError.WithMessage(err.Error())
		}
	} else {
		err := m.Client.RemoveObject(ctx, payload.BucketName, payload.Dir, minio.RemoveObjectOptions{})
		if err != nil {
			logger.Error.Printf("Error while deleting object %s in bucket %s: %v", payload.Dir, payload.BucketName, err)
			return errors.ExternalServiceError.WithMessage(err.Error())
		}
	}
	return errors.OK
}

func (m *Minio) Close(_ context.Context) errors.Error {
	logger.Info.Println("Nothing to do")
	return errors.OK
}
