package gcp

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
)

type gcpRepository struct {
	sshKeysRepository repository.SSHKeys
}

func NewRepository(sshKeysRepo repository.SSHKeys) resourceRepo.Resource {
	return &gcpRepository{}
}
