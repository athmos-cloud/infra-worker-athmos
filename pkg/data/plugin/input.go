package plugin

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"

const (
	managedKey    = "managed"
	tagsKey       = "tags"
	identifierKey = "identifier"
)

type BasePayload struct {
	Identifier identifier.ID
	Managed    bool
	Tags       map[string]string
}

func InputWithBaseFields(input *map[string]interface{}, payload BasePayload) {
	inputValue := *input
	inputValue[managedKey] = payload.Managed
	inputValue[identifierKey] = payload.Identifier
	if payload.Tags != nil {
		inputValue[tagsKey] = payload.Tags
	} else {
		inputValue[tagsKey] = make(map[string]string)
	}
	*input = inputValue
}
