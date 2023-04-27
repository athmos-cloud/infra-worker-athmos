package kubernetes

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
	"strings"
)

const (
	manifestSeparatorString = "---"
)

type Identifier struct {
	ResourceID schema.GroupVersionResource
	Name       string
	Namespace  string
}

func (identifier Identifier) Equals(other Identifier) bool {
	return identifier.ResourceID == other.ResourceID &&
		identifier.Name == other.Name &&
		identifier.Namespace == other.Namespace
}

func GetResourcesIdentifiersFromManifests(manifests string) ([]Identifier, errors.Error) {
	var identifierList []Identifier
	manifestList := strings.Split(manifests, manifestSeparatorString)
	for _, val := range manifestList {
		identifier, err := getResourceIdentifierFromManifest(val)
		if !err.IsOk() {
			return identifierList, err
		}
		identifierList = append(identifierList, identifier)
	}
	return identifierList, errors.OK
}

func getResourceIdentifierFromManifest(manifest string) (Identifier, errors.Error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, gvk, err := decode([]byte(manifest), nil, nil)
	if err != nil {
		log.Fatal("Error decoding manifest:", err)
		return Identifier{}, errors.InvalidArgument.WithMessage(err.Error())
	}
	groupVersionResource := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: gvk.Kind, // strings.ToLower(gvk.Kind) + "s" ??
	}

	unstructuredObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		log.Fatalf("Error casting to unstructured.Unstructured")
	}
	name := unstructuredObj.GetName()
	namespace := unstructuredObj.GetNamespace()

	return Identifier{
		ResourceID: groupVersionResource,
		Name:       name,
		Namespace:  namespace,
	}, errors.OK

}
