package dto

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
)

type GetProviderRequest struct {
	IdentifierID identifier.Provider `json:"identifierID"`
}

type GetProviderResponse struct {
	ProjectID string            `json:"projectID"`
	Payload   resource.Provider `json:"payload"`
}

type ListProvidersRequest struct {
	Recursive bool `json:"recursive" default:"false"`
}

type ListProvidersResponse struct {
	ProjectID string                      `json:"projectID"`
	Payload   resource.ProviderCollection `json:"payload"`
}

type GetProviderStackRequest struct {
	ProviderID string `json:"providerID"`
}

type CreateProviderRequest struct {
	Name           string `json:"name"`
	VPC            string `json:"vpc,omitempty"`
	SecretAuthName string `json:"secretAuthName"`
}

type CreateProviderResponse struct {
	ProjectID string            `json:"projectID"`
	Payload   resource.Provider `json:"payload"`
}

type UpdateProviderRequest struct {
	IdentifierID   identifier.Provider `json:"identifierID"`
	Name           string              `json:"name"`
	SecretAuthName string              `json:"secretAuthName"`
}

type DeleteProviderRequest struct {
	IdentifierID identifier.Provider `json:"identifierID"`
	Cascade      bool                `json:"cascade" default:"false"`
}
