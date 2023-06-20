package aws

import (
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
)

type awsRepository struct{}


func NewRepository() resourceRepo.Resource {
	return &awsRepository{}
}
