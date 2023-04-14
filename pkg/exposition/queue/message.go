package queue

type Verb string

const (
	CREATE Verb = "create"
	UPDATE Verb = "update"
	DELETE Verb = "delete"
)

type Message struct {
	Verb    Verb        `json:"verb"`
	Payload interface{} `json:"payload"`
}
