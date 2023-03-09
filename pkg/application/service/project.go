package service

import (
	"github.com/PaulBarrie/infra-worker/pkg/application/service/dto"
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/runtime"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
)

type ProjectService struct {
	Runtime runtime.IRuntime
}

func NewProjectService(runtime runtime.IRuntime) *ProjectService {
	return &ProjectService{
		Runtime: runtime,
	}
}

type ProjectServiceRequestPayload struct {
	ProjectId string
}

func (p *ProjectService) Add(payload dto.CreatePluginInstanceRequest) errors.Error {
	//Retrieve the project
	// Retrieve the plugin

	// Build the plugin -> pipeline
	return errors.OK
}
