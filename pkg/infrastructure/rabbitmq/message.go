package rabbitmq

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
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
	Identifier   identifier.Payload  `json:"identifier,omitempty"`
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
