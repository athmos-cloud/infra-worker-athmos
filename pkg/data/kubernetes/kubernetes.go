package kubernetes

type Resource struct {
	Identifier Identifier
	Events     EventList
	Outputs    OutputList
}

func NewResourceList(identifiers []Identifier) ResourceList {
	resources := make(ResourceList, len(identifiers))
	for i, identifier := range identifiers {
		resources[i] = Resource{Identifier: identifier}
	}
	return resources
}

func (resource *Resource) Equals(other Resource) bool {
	return resource.Identifier.Equals(other.Identifier) &&
		resource.Events.Equals(other.Events) &&
		resource.Outputs.Equals(other.Outputs)
}

type ResourceList []Resource

func (resourceList *ResourceList) Equals(other ResourceList) bool {
	if len(*resourceList) != len(other) {
		return false
	}
	for i, resource := range *resourceList {
		if !resource.Equals(other[i]) {
			return false
		}
	}
	return true
}
