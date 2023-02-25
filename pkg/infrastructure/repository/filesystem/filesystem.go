package filesystem

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"os"
	"path/filepath"
)

type Filesystem struct {
}
type CreatePayload struct {
	DestPath string
	Content  string
}

func (f Filesystem) Get() repository.IRepository {
	//TODO implement me
	panic("implement me")
}

func (f Filesystem) Create(ctx context.Context, options option.List) errors.Error {
	if !options.Validate(2) {
		return errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"expects list of destPath and content str. Got : %v", options,
			),
		)
	}
	destPath := options.Values[0].(string)
	content := []byte(options.Values[1].(string))

	folder := filepath.Base(destPath)
	if err := os.MkdirAll(folder, 0777); err != nil {
		return errors.IOError.WithMessage(err)
	}
	file, err := os.Create(destPath)
	if err != nil {
		return errors.IOError.WithMessage(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			errors.IOError.WithMessage(err)
		}
	}(file)

	if err = gob.NewEncoder(file).Encode(content); err != nil {
		return errors.IOError.WithMessage(err)
	}

	return errors.OK
}

func (f Filesystem) GetObject(options option.List) (*[]byte, errors.Error) {
	if !options.Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"expects list of filePath str. Got : %v", options,
			),
		)
	}
	filePath := options.Values[0].(string)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.IOError.WithMessage(err)
	}
	return &content, errors.OK

}

func (f Filesystem) Update(ctx context.Context, s ...string) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (f Filesystem) Delete(ctx context.Context, s ...string) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (f Filesystem) Close() errors.Error {
	//TODO implement me
	panic("implement me")
}
