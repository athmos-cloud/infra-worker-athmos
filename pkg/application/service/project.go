package service

import (
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
)

type ProjectService struct {
}

type ProjectServiceRequestPayload struct {
	ProjectId string
}

func (p *ProjectService) Add() errors.Error {
	//Retrieve the project
	// Retrieve the plugin

	// Build the plugin -> pipeline
	return errors.OK
}
