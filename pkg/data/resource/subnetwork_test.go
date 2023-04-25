package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
	"testing"
)

func TestNewSubnetwork(t *testing.T) {
	type args struct {
		id identifier.Subnetwork
	}
	tests := []struct {
		name string
		args args
		want Subnetwork
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSubnetwork(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSubnetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubnetwork_FromMap(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Subnetwork
		KubernetesResources kubernetes.ResourceList
		VPC                 string
		Network             string
		Region              string
		IPCIDRRange         string
		VMs                 VMCollection
	}
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subnet := &Subnetwork{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				VPC:                 tt.fields.VPC,
				Network:             tt.fields.Network,
				Region:              tt.fields.Region,
				IPCIDRRange:         tt.fields.IPCIDRRange,
				VMs:                 tt.fields.VMs,
			}
			if got := subnet.FromMap(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubnetwork_GetMetadata(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Subnetwork
		KubernetesResources kubernetes.ResourceList
		VPC                 string
		Network             string
		Region              string
		IPCIDRRange         string
		VMs                 VMCollection
	}
	tests := []struct {
		name   string
		fields fields
		want   metadata.Metadata
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subnet := &Subnetwork{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				VPC:                 tt.fields.VPC,
				Network:             tt.fields.Network,
				Region:              tt.fields.Region,
				IPCIDRRange:         tt.fields.IPCIDRRange,
				VMs:                 tt.fields.VMs,
			}
			if got := subnet.GetMetadata(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubnetwork_GetPluginReference(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Subnetwork
		KubernetesResources kubernetes.ResourceList
		VPC                 string
		Network             string
		Region              string
		IPCIDRRange         string
		VMs                 VMCollection
	}
	type args struct {
		request resource.GetPluginReferenceRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   resource.GetPluginReferenceResponse
		want1  errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subnet := &Subnetwork{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				VPC:                 tt.fields.VPC,
				Network:             tt.fields.Network,
				Region:              tt.fields.Region,
				IPCIDRRange:         tt.fields.IPCIDRRange,
				VMs:                 tt.fields.VMs,
			}
			got, got1 := subnet.GetPluginReference(tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPluginReference() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetPluginReference() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSubnetwork_Insert(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Subnetwork
		KubernetesResources kubernetes.ResourceList
		VPC                 string
		Network             string
		Region              string
		IPCIDRRange         string
		VMs                 VMCollection
	}
	type args struct {
		project Project
		update  []bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subnet := &Subnetwork{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				VPC:                 tt.fields.VPC,
				Network:             tt.fields.Network,
				Region:              tt.fields.Region,
				IPCIDRRange:         tt.fields.IPCIDRRange,
				VMs:                 tt.fields.VMs,
			}
			if got := subnet.Insert(tt.args.project, tt.args.update...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubnetwork_ToDomain(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Subnetwork
		KubernetesResources kubernetes.ResourceList
		VPC                 string
		Network             string
		Region              string
		IPCIDRRange         string
		VMs                 VMCollection
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
		want1  errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subnet := &Subnetwork{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				VPC:                 tt.fields.VPC,
				Network:             tt.fields.Network,
				Region:              tt.fields.Region,
				IPCIDRRange:         tt.fields.IPCIDRRange,
				VMs:                 tt.fields.VMs,
			}
			got, got1 := subnet.ToDomain()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDomain() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ToDomain() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSubnetwork_WithMetadata(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Subnetwork
		KubernetesResources kubernetes.ResourceList
		VPC                 string
		Network             string
		Region              string
		IPCIDRRange         string
		VMs                 VMCollection
	}
	type args struct {
		request metadata.CreateMetadataRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subnet := &Subnetwork{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				VPC:                 tt.fields.VPC,
				Network:             tt.fields.Network,
				Region:              tt.fields.Region,
				IPCIDRRange:         tt.fields.IPCIDRRange,
				VMs:                 tt.fields.VMs,
			}
			subnet.WithMetadata(tt.args.request)
		})
	}
}
