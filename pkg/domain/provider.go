package domain

import "github.com/athmos-cloud/infra-worker-athmos/pkg/common"

type Provider struct {
	ID           string
	Name         string
	ProviderType common.ProviderType
	Projects     map[string]Project
}
