package rabbitmq

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"time"
)

const (
	nestPatternValue = "project-event"
)

type Verb string

const (
	CREATE Verb = "create"
	UPDATE Verb = "update"
	DELETE Verb = "delete"
)

type messageReceived struct {
	Pattern pattern     `json:"pattern"`
	Data    dataMessage `json:"data"`
}

type pattern struct {
	Pattern string `json:"pattern"`
}

type dataMessage struct {
	Verb         Verb   `json:"verb"`
	ProjectID    string `json:"project_id"`
	ProviderType string `json:"provider_type"`
	ResourceType string `json:"resource_type"`
	Payload      any    `json:"payload"`
}

type MessageSend struct {
	ProjectID    string              `json:"project_id"`
	ResourceType types.Resource      `json:"resource_type"`
	ProviderType types.Provider      `json:"provider_type"`
	Code         int                 `json:"code"`
	Type         metadata.StatusType `json:"type"`
	Date         time.Time           `json:"date"`
	Message      string              `json:"message"`
	Identifier   identifier.ID       `json:"identifier,omitempty"`
	Payload      interface{}         `json:"payload,omitempty"`
}

type NestMessageWrap struct {
	Pattern string `json:"pattern"`
	Data    any    `json:"data"`
}

func (ms *MessageSend) WithNestWrapper() NestMessageWrap {
	return NestMessageWrap{
		Pattern: nestPatternValue,
		Data:    ms,
	}
}

func getResourceID(resourcePayload any, resourceType types.Resource) identifier.ID {
	switch resourceType {
	case types.ProviderResource:
		provider := resourcePayload.(resource.Provider)
		return &provider.IdentifierID
	case types.NetworkResource:
		network := resourcePayload.(network.Network)
		return &network.IdentifierID
	case types.FirewallResource:
		firewall := resourcePayload.(network.Firewall)
		return &firewall.IdentifierID
	case types.SubnetworkResource:
		subnetwork := resourcePayload.(network.Subnetwork)
		return &subnetwork.IdentifierID
	case types.VMResource:
		vm := resourcePayload.(instance.VM)
		return &vm.IdentifierID
	case types.SqlDBResource:
		sqldb := resourcePayload.(instance.SqlDB)
		return &sqldb.IdentifierID
	}

	return nil
}

func getResourceName(resourcePayload any, resourceType types.Resource) identifier.ID {
	switch resourceType {
	case types.ProviderResource:
		provider := resourcePayload.(resource.Provider)
		return &provider.IdentifierName
	case types.NetworkResource:
		network := resourcePayload.(network.Network)
		return &network.IdentifierName
	case types.FirewallResource:
		firewall := resourcePayload.(network.Firewall)
		return &firewall.IdentifierName
	case types.SubnetworkResource:
		subnetwork := resourcePayload.(network.Subnetwork)
		return &subnetwork.IdentifierName
	case types.VMResource:
		vm := resourcePayload.(instance.VM)
		return &vm.IdentifierName
	case types.SqlDBResource:
		sqldb := resourcePayload.(instance.SqlDB)
		return &sqldb.IdentifierName
	}

	return nil
}
