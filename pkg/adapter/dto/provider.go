package dto

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
)

type GetProviderRequest struct {
	ProjectID  string              `json:"projectID"`
	Identifier identifier.Provider `json:"identifier"`
}

type GetProviderResponse struct {
	ProjectID string            `json:"projectID"`
	Payload   resource.Provider `json:"payload"`
}

type GetAllProvidersRequest struct {
	ProjectID string `json:"projectID"`
	Recursive bool   `json:"recursive"`
}

type GetAllProvidersResponse struct {
	ProjectID string                      `json:"projectID"`
	Payload   resource.ProviderCollection `json:"payload"`
}

type CreateProviderRequest struct {
	ProjectID      string `json:"projectID"`
	Name           string `json:"name"`
	SecretAuthName string `json:"secretAuthName"`
}

type CreateProviderResponse struct {
	ProjectID string            `json:"projectID"`
	Payload   resource.Provider `json:"payload"`
}

type UpdateProviderRequest struct {
	ProjectID      string              `json:"projectID"`
	Identifier     identifier.Provider `json:"identifier"`
	Name           string              `json:"name"`
	SecretAuthName string              `json:"secretAuthName"`
}

type DeleteProviderRequest struct {
	ProjectID  string              `json:"projectID"`
	Identifier identifier.Provider `json:"identifier"`
	Cascade    bool                `json:"cascade"`
}
