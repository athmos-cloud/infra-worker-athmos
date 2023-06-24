package rabbitmq

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
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
	Pattern string      `json:"pattern"`
	Data    dataMessage `json:"data"`
}

type dataMessage struct {
	Verb         Verb   `json:"verb"`
	ProjectID    string `json:"project_id"`
	ProviderType string `json:"provider_type"`
	ResourceType string `json:"resource_type"`
	Payload      any    `json:"payload"`
}

type eventType string

const (
	EventTypeCreateRequestSent eventType = "CREATE_REQUEST_SENT"
	EventTypeUpdateRequestSent eventType = "UPDATE_REQUEST_SENT"
	EventTypeDeleteRequestSent eventType = "DELETE_REQUEST_SENT"

	// Error CreateRequestTreated eventType = "CREATE_REQUEST_TREATED"
	//ResourceCreated      eventType = "RESOURCE_CREATED"
	Error eventType = "ERROR"
)

type MessageSend struct {
	ProjectID  string             `json:"project_id"`
	Code       int                `json:"code"`
	Type       eventType          `json:"type"`
	Date       time.Time          `json:"date"`
	Message    string             `json:"message"`
	Identifier identifier.Payload `json:"identifier,omitempty"`
	Payload    interface{}        `json:"payload,omitempty"`
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
