package rabbitmq

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

const (
	projectIDKey        = "projectID"
	nameKey             = "name"
	managedKey          = "managed"
	tagsKey             = "tags"
	parentIdentifierKey = "parentIdentifier"
	providerTypeKey     = "providerType"
	resourceTypeKey     = "resourceType"
	resourceSpecsKey    = "resourceSpecs"
)

func mapToCreateResourceRequest(entry map[string]interface{}) resource.CreateResourceRequest {
	projectID := ""
	name := ""
	managed := true
	tags := map[string]string{}
	var parentIdentifier identifier.ID
	var providerType types.ProviderType
	var resourceType types.ResourceType
	resourceSpecs := map[string]interface{}{}

	if _, ok := entry[projectIDKey]; ok {
		projectID = entry[projectIDKey].(string)
	}
	if _, ok := entry[nameKey]; ok {
		name = entry[nameKey].(string)
	}
	if _, ok := entry[managedKey]; ok {
		if managed, ok = entry[managedKey].(bool); !ok {
			panic(errors.BadRequest.WithMessage("managed must be a boolean"))
		}
	}
	if _, ok := entry[tagsKey]; ok {
		if _, okCast := entry[tagsKey].(map[string]string); okCast {
			tags = entry[tagsKey].(map[string]string)
		}
	}
	if _, ok := entry[parentIdentifierKey]; ok {
		if _, okCast := entry[parentIdentifierKey].(map[string]interface{}); okCast {
			parentIdentifier = identifier.FromMap(entry[parentIdentifierKey].(map[string]interface{}))
		} else {
			panic(errors.BadRequest.WithMessage("parentIdentifier must be a map[string]interface{}"))
		}
	} else {
		parentIdentifier = identifier.Empty{}
	}
	if _, ok := entry[providerTypeKey]; ok {
		providerTypeEntry, okCast := entry[providerTypeKey].(string)
		if !okCast {
			panic(errors.BadRequest.WithMessage("providerType must be a string"))
		}
		providerType = types.ProviderTypeFromString(providerTypeEntry)
	}
	if _, ok := entry[resourceTypeKey]; ok {
		providerTypeEntry, okCast := entry[resourceTypeKey].(string)
		if !okCast {
			panic(errors.BadRequest.WithMessage("providerType must be a string"))
		}
		resourceType = types.ResourceTypeFromString(providerTypeEntry)
	}
	if _, ok := entry[resourceSpecsKey]; ok {
		if _, okCast := entry[resourceSpecsKey].(map[string]interface{}); okCast {
			resourceSpecs = entry[resourceSpecsKey].(map[string]interface{})
		}
	}

	return resource.CreateResourceRequest{
		ProjectID:        projectID,
		Name:             name,
		Managed:          managed,
		Tags:             tags,
		ParentIdentifier: parentIdentifier,
		ProviderType:     providerType,
		ResourceType:     resourceType,
		ResourceSpecs:    resourceSpecs,
	}
}
