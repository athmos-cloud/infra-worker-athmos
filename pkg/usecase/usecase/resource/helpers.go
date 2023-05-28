package resourceUc

import (
	"fmt"
	"math/rand"
)

func idFromName(name string) string {
	return fmt.Sprintf("%s-%s", name, randomString(8))
}

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter)-1)]
	}
	return string(b)
}