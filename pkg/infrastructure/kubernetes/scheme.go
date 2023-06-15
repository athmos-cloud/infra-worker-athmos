package kubernetes

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	computeGCP "github.com/upbound/provider-gcp/apis/compute/v1beta1"
	sqlGCP "github.com/upbound/provider-gcp/apis/sql/v1beta1"
	providerGCP "github.com/upbound/provider-gcp/apis/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

func getScheme() *runtime.Scheme {
	newScheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(newScheme))
	registerGCPResources(newScheme)

	return newScheme
}

func registerGCPResources(runtimeScheme *runtime.Scheme) {
	UpboundSchemeBuilder := &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "gcp.upbound.io", Version: "v1beta1"}}
	UpboundSchemeBuilder.Register(&providerGCP.ProviderConfig{}, &providerGCP.ProviderConfigList{})
	if err := UpboundSchemeBuilder.AddToScheme(runtimeScheme); err != nil {
		panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error registering GCP upbound resources: %v", err)))
	}
	UpboundComputeSchemeBuilder := &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "compute.gcp.upbound.io", Version: "v1beta1"}}
	UpboundComputeSchemeBuilder.Register(&computeGCP.Network{}, &computeGCP.NetworkList{})
	UpboundComputeSchemeBuilder.Register(&computeGCP.Subnetwork{}, &computeGCP.SubnetworkList{})
	UpboundComputeSchemeBuilder.Register(&computeGCP.Firewall{}, &computeGCP.FirewallList{})
	UpboundComputeSchemeBuilder.Register(&computeGCP.Instance{}, &computeGCP.InstanceList{})
	if err := UpboundComputeSchemeBuilder.AddToScheme(runtimeScheme); err != nil {
		panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error registering GCP compute resources: %v", err)))
	}
	UpboundSqlSchemeBuilder := &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "sql.gcp.upbound.io", Version: "v1beta1"}}
	UpboundSqlSchemeBuilder.Register(&sqlGCP.DatabaseInstance{}, &sqlGCP.DatabaseInstanceList{})
	if err := UpboundSqlSchemeBuilder.AddToScheme(runtimeScheme); err != nil {
		panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error registering GCP sql resources: %v", err)))
	}
}
