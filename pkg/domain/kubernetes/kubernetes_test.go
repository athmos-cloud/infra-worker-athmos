package kubernetes

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

func TestGetResourcesIdentifiersFromManifests(t *testing.T) {
	manifests := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example-deployment
  namespace: example-namespace
---
apiVersion: v1
kind: Service
metadata:
  name: example-service
  namespace: example-namespace
`

	expected := []Identifier{
		{
			ResourceID: schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "Deployment",
			},
			Name:      "example-deployment",
			Namespace: "example-namespace",
		},
		{
			ResourceID: schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "Service",
			},
			Name:      "example-service",
			Namespace: "example-namespace",
		},
	}

	identifiers, err := GetResourcesIdentifiersFromManifests(manifests)

	if !err.IsOk() {
		t.Errorf("GetResourcesIdentifiersFromManifests returned an error: %v", err)
	}

	if len(identifiers) != len(expected) {
		t.Errorf("Expected %d identifiers, got %d", len(expected), len(identifiers))
	}

	for i := range expected {
		if identifiers[i].ResourceID != expected[i].ResourceID ||
			identifiers[i].Name != expected[i].Name ||
			identifiers[i].Namespace != expected[i].Namespace {
			t.Errorf("Identifier mismatch. Expected: %+v, Got: %+v", expected[i], identifiers[i])
		}
	}
}
