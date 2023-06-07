package rabbitmq

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
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
	Pattern string      `json:"pattern"`
	Data    dataMessage `json:"data"`
}

type dataMessage struct {
	Verb         Verb           `json:"verb"`
	ProjectID    string         `json:"project_id"`
	ProviderType types.Provider `json:"provider_type"`
	ResourceType types.Resource `json:"resource_type"`
	Payload      any            `json:"payload"`
}
type eventType string

const (
	eventTypeCreateRequestSent eventType = "CREATE_REQUEST_SENT"
	eventTypeUpdateRequestSent eventType = "UPDATE_REQUEST_SENT"
	eventTypeDeleteRequestSent eventType = "DELETE_REQUEST_SENT"

	// Error CreateRequestTreated eventType = "CREATE_REQUEST_TREATED"
	//ResourceCreated      eventType = "RESOURCE_CREATED"
	Error eventType = "ERROR"
)

type messageSend struct {
	ProjectID  string             `json:"project_id"`
	Code       int                `json:"code"`
	Type       eventType          `json:"type"`
	Identifier identifier.Payload `json:"identifier,omitempty"`
	Payload    interface{}        `json:"payload,omitempty"`
}

type nestMessageWrap struct {
	Pattern string `json:"pattern"`
	Data    any    `json:"data"`
}

func (ms *messageSend) WithNestWrapper() nestMessageWrap {
	return nestMessageWrap{
		Pattern: nestPatternValue,
		Data:    ms,
	}
}
