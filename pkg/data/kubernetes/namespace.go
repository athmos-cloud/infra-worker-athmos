package kubernetes

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"regexp"
)

func NamespaceFormat(namespace string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(errors.InternalError.WithMessage(err.Error()))
	}
	return reg.ReplaceAllString(namespace, "")
}
