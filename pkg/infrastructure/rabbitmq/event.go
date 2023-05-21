package rabbitmq

type EventType string

const (
	CreateRequestSent    EventType = "CREATE_REQUEST_SENT"
	CreateRequestTreated EventType = "CREATE_REQUEST_TREATED"
	ResourceCreated      EventType = "RESOURCE_CREATED"
	CreateError          EventType = "CREATE_ERROR"
	Error                EventType = "ERROR"
)

type Event struct {
	ProjectID string      `json:"project_id"`
	Code    int         `json:"code"`
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}
