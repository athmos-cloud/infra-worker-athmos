package types

import (
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"strings"
)

type ProviderType string

const (
	AWS   ProviderType = "aws"
	Azure ProviderType = "azure"
	GCP   ProviderType = "gcp"
)

func (p ProviderType) FromString(str string) ProviderType {
	strShaped := strings.ReplaceAll(strings.ToLower(str), " ", "")
	switch strShaped {
	case string(AWS):
		return AWS
	case string(Azure):
		return Azure
	case string(GCP):
		return GCP
	default:
		logger.Warning.Printf("Provider type %s not recognised, defaulting to AWS", str)
		return ""
	}
}
