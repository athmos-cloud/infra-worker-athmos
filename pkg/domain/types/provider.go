package types

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type Provider string

const (
	ProviderAWS   Provider = "aws"
	ProviderGCP   Provider = "gcp"
	ProviderAZURE Provider = "azure"
)

var providersMapping = map[string]Provider{
	"aws":   ProviderAWS,
	"gcp":   ProviderGCP,
	"azure": ProviderAZURE,
}

func ProviderFromString(s string) (Provider, errors.Error) {
	if val, ok := providersMapping[s]; ok {
		return val, errors.OK
	}
	return "", errors.BadRequest.WithMessage(fmt.Sprintf("provider %s is not supported", s))
}
