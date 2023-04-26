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
		subnetwork Subnetwork
	}
	type args struct {
		project Project
		update  []bool
	}
	type want struct {
		err        errors.Error
		subnetwork Subnetwork
	}
	providerID := "test"
	vpcID := "test"
	networkID := "test"

	subnetwork1 := NewSubnetwork(identifier.Subnetwork{ID: "test-1", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID})
	subnetwork2 := NewSubnetwork(identifier.Subnetwork{ID: "test-2", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID})
	subnetwork3 := subnetwork1
	subnetwork3.Metadata.Tags = map[string]string{"test": "test"}
	subnetwork4 := subnetwork3
	subnetwork4.Metadata.Tags = map[string]string{"hello": "world"}
	subnetwork5 := NewSubnetwork(identifier.Subnetwork{ID: "test-5", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID})

	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(identifier.Provider{ID: providerID})
	testVPC := NewVPC(identifier.VPC{ID: vpcID, ProviderID: providerID})
	testNetwork := NewNetwork(identifier.Network{ID: networkID, ProviderID: providerID, VPCID: vpcID})
	testNetwork.Subnetworks[subnetwork1.Identifier.ID] = subnetwork1
	testVPC.Networks[networkID] = testNetwork
	testProvider.VPCs[vpcID] = testVPC
	testProject.Resources[providerID] = testProvider

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Insert non existing subnetwork (creation)",
			fields: fields{
				subnetwork: subnetwork2,
			},
			args: args{
				testProject,
				[]bool{},
			},
			want: want{
				err:        errors.OK,
				subnetwork: subnetwork2,
			},
		},
		{
			name: "Update existing subnetwork (update)",
			fields: fields{
				subnetwork: subnetwork3,
			},
			args: args{
				testProject,
				[]bool{true},
			},
			want: want{
				err:        errors.OK,
				subnetwork: subnetwork3,
			},
		},
		{
			name: "Update existing subnetwork (no update)",
			fields: fields{
				subnetwork: subnetwork4,
			},
			args: args{
				testProject,
				[]bool{},
			},
			want: want{
				err:        errors.Conflict,
				subnetwork: subnetwork3,
			},
		},
		{
			name: "Update non existing subnetwork",
			fields: fields{
				subnetwork: subnetwork5,
			},
			args: args{
				testProject,
				[]bool{true},
			},
			want: want{
				err:        errors.NotFound,
				subnetwork: Subnetwork{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subnet := tt.fields.subnetwork
			if got := subnet.Insert(tt.args.project, tt.args.update...); !reflect.DeepEqual(got.Code, tt.want.err.Code) {
				t.Errorf("Insert() = %v, want %v", got.Code, tt.want.err.Code)
			}
			id := tt.fields.subnetwork.Identifier
			if !reflect.DeepEqual(testProject.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.ID], tt.want.subnetwork) {
				t.Errorf("Insert() = %v, want %v", subnet, tt.want.subnetwork)
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
