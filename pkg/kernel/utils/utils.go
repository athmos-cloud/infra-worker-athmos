package utils

import (
	"github.com/google/uuid"
	"math/rand"
)

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateUUID() string {
	return uuid.New().String()
}
