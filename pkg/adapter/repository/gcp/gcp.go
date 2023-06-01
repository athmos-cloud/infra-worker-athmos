package gcp

import (
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
)

type gcpRepository struct{}

func NewRepository() resourceRepo.Resource {
	return &gcpRepository{}
}
