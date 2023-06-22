package usecase

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
)

func IdFromName(name string) string {
	return fmt.Sprintf("%s-%s", strings.ToLower(RemoveSpecialChars(name)), RandomString(8))
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter)-1)]
	}
	return string(b)
}

func RemoveSpecialChars(s string) string {
	// The regular expression matches any character that is not a letter or number.
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(s, "")
	return processedString
}
