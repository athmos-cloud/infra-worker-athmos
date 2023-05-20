package gcp

import (
	_ "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
)

type gcpRepository struct{}

func NewRepository() resourceRepo.ProviderResource {
	return &gcpRepository{}
}
