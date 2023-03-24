package utils

import (
	"fmt"
	config2 "github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"math/rand"
	"os"
)

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func StringToTempFile(content string) (*os.File, errors.Error) {
	file, err := os.CreateTemp(fmt.Sprintf("%s/%s", config2.Get().TempDir, RandomString(10)), "temp")
	if err != nil {
		return nil, errors.IOError.WithMessage(err)
	}

	_, err = file.WriteString(content)
	if err != nil {
		return nil, errors.IOError.WithMessage(err)

	}
	return file, errors.OK
}
