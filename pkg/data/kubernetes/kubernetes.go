package kubernetes

type ResourceList []Resource

type Resource struct {
	Identifier Identifier
	Events     EventList
	Outputs    OutputList
}
