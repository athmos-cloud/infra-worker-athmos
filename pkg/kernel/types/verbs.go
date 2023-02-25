package types

import (
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"strings"
)

type Verbs string

const (
	Get    Verbs = "get"
	Create Verbs = "create"
	Update Verbs = "update"
	Delete Verbs = "delete"
)

func (v Verbs) FromString(str string) Verbs {
	strShaped := strings.ReplaceAll(strings.ToLower(str), " ", "")
	switch strShaped {
	case string(Get):
		return Get
	case string(Create):
		return Create
	case string(Update):
		return Update
	case string(Delete):
		return Delete
	default:
		logger.Warning.Printf("Verb %s not recognised, defaulting to AWS", str)
		return ""
	}
}
