package utils

import (
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

func MapEquals(a map[string]string, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if b[key] != value {
			return false
		}
	}
	return true
}

func SliceEquals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, value := range a {
		equals := false
		for _, otherValue := range b {
			if value == otherValue {
				equals = true
			}
		}
		if !equals {
			return false
		}
	}
	return true
}

func IntSliceEquals(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for _, value := range a {
		equals := false
		for _, otherValue := range b {
			if value == otherValue {
				equals = true
			}
		}
		if !equals {
			return false
		}
	}
	return true
}

func MapSafeInsert(m *map[string]interface{}, key string, value interface{}) {
	curMap := *m
	if curMap == nil {
		curMap = make(map[string]interface{})
	}
	curMap[key] = value
	*m = curMap
}
