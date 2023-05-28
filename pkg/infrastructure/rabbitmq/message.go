package rabbitmq

import "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"

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
	Verb Verb        `json:"verb"`
	Data dataMessage `json:"payload"`
}

type dataMessage struct {
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
	ProjectID string      `json:"project_id"`
	Code      int         `json:"code"`
	Type      eventType   `json:"type"`
	Payload   interface{} `json:"payload"`
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
