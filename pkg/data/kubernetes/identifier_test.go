package kubernetes

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

func TestGetResourcesIdentifiersFromManifests(t *testing.T) {
	manifests := `
apiVersion: gcp.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  annotations:
    meta.helm.sh/release-name: gcp-provider-pkvvd
    meta.helm.sh/release-namespace: gcptesttedwk
  creationTimestamp: "2023-04-30T17:00:50Z"
  finalizers:
  - in-use.crossplane.io
  generation: 1
  labels:
    app.kubernetes.io/managed-by: Helm
  name: gcp-provider-1
  namespace: test1
  resourceVersion: "900102"
  uid: 47eedb8d-9762-43a8-9e3b-28f620d4da43
spec:
  credentials:
    secretRef:
      key: credentials.json
      name: gcp-secret
      namespace: gcptesttedwk
    source: Secret
  projectID: plugin-playground
status: {}
---
apiVersion: gcp.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  annotations:
    meta.helm.sh/release-name: gcp-provider-pkvvd
    meta.helm.sh/release-namespace: gcptesttedwk
  creationTimestamp: "2023-04-30T17:00:50Z"
  finalizers:
  - in-use.crossplane.io
  generation: 1
  labels:
    app.kubernetes.io/managed-by: Helm
  name: gcp-provider-2
  namespace: test2
  resourceVersion: "900102"
  uid: 47eedb8d-9762-43a8-9e3b-28f620d4da43
spec:
  credentials:
    secretRef:
      key: credentials.json
      name: gcp-secret
      namespace: gcptesttedwk
    source: Secret
  projectID: plugin-playground
status: {}


`

	expected := []Identifier{
		{
			ResourceID: schema.GroupVersionResource{
				Group:    "gcp.upbound.io",
				Version:  "v1beta1",
				Resource: "ProviderConfig",
			},
			Name:      "gcp-provider-1",
			Namespace: "test1",
		},
		{
			ResourceID: schema.GroupVersionResource{
				Group:    "gcp.upbound.io",
				Version:  "v1beta1",
				Resource: "ProviderConfig",
			},
			Name:      "gcp-provider-2",
			Namespace: "test2",
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: %v", r)
		}
	}()
	identifiers := GetResourcesIdentifiersFromManifests(manifests)

	if len(identifiers) != len(expected) {
		t.Errorf("Expected %d identifiers, got %d", len(expected), len(identifiers))
	}

	for i := range expected {
		if identifiers[i].ResourceID != expected[i].ResourceID ||
			identifiers[i].Name != expected[i].Name ||
			identifiers[i].Namespace != expected[i].Namespace {
			t.Errorf("ParentIdentifier mismatch. Expected: %+v, Got: %+v", expected[i], identifiers[i])
		}
	}
}
