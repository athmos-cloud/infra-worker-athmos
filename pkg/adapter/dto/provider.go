package dto

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
)

type GetProviderResponse struct {
	ProjectID string            `json:"project_id"`
	Payload   resource.Provider `json:"payload"`
}

type ListProvidersRequest struct{}

type ListProvidersResponse struct {
	ProjectID string                      `json:"project_id"`
	Payload   resource.ProviderCollection `json:"payload"`
}

type ListProviderItemResponse struct {
	Name       string              `json:"name"`
	Identifier identifier.Provider `json:"identifier"`
}

type GetProviderStackRequest struct {
	ProviderID string `json:"provider_id"`
}

type GetProviderStackResponse struct {
	ProjectID string            `json:"project_id"`
	Payload   resource.Provider `json:"payload"`
}

type CreateProviderRequest struct {
	Name           string `json:"name"`
	VPC            string `json:"vpc,omitempty"`
	SecretAuthName string `json:"secret_auth_name"`
}

type CreateProviderResponse struct {
	ProjectID string            `json:"project_id"`
	Payload   resource.Provider `json:"payload"`
}

type UpdateProviderRequest struct {
	IdentifierID   identifier.Provider `json:"identifier_id"`
	Name           string              `json:"name"`
	SecretAuthName string              `json:"secret_auth_name"`
}

type DeleteProviderRequest struct {
	IdentifierID identifier.Provider `json:"identifier_id"`
	Cascade      bool                `json:"cascade" default:"false"`
}
