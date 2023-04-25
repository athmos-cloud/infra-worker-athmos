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

func TestNetwork_FromMap(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Network
		KubernetesResources kubernetes.ResourceList
		Subnetworks         SubnetworkCollection
		Firewalls           FirewallCollection
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
			network := &Network{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Subnetworks:         tt.fields.Subnetworks,
				Firewalls:           tt.fields.Firewalls,
			}
			if got := network.FromMap(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetwork_GetMetadata(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Network
		KubernetesResources kubernetes.ResourceList
		Subnetworks         SubnetworkCollection
		Firewalls           FirewallCollection
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
			network := &Network{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Subnetworks:         tt.fields.Subnetworks,
				Firewalls:           tt.fields.Firewalls,
			}
			if got := network.GetMetadata(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetwork_GetPluginReference(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Network
		KubernetesResources kubernetes.ResourceList
		Subnetworks         SubnetworkCollection
		Firewalls           FirewallCollection
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
			network := &Network{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Subnetworks:         tt.fields.Subnetworks,
				Firewalls:           tt.fields.Firewalls,
			}
			got, got1 := network.GetPluginReference(tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPluginReference() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetPluginReference() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetwork_Insert(t *testing.T) {
	type fields struct {
		network Network
	}
	type args struct {
		project Project
		update  []bool
	}
	type want struct {
		err     errors.Error
		network Network
	}

	providerID := "test"
	vpcID := "test"

	network1 := NewNetwork(identifier.Network{ID: "test-1", ProviderID: providerID, VPCID: vpcID})
	network2 := NewNetwork(identifier.Network{ID: "test-2", ProviderID: providerID, VPCID: vpcID})
	network3 := network1
	network3.Metadata.Tags = map[string]string{"test": "test"}
	network4 := network3
	network4.Metadata.Tags = map[string]string{"hello": "world"}
	network5 := NewNetwork(identifier.Network{ID: "test-5", ProviderID: providerID, VPCID: vpcID})

	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(identifier.Provider{ID: providerID})
	testVPC := NewVPC(identifier.VPC{ID: vpcID, ProviderID: providerID})
	testVPC.Networks[network1.Identifier.ID] = network1
	testProvider.VPCs[vpcID] = testVPC
	testProject.Resources[providerID] = testProvider

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Insert non existing network (creation)",
			fields: fields{
				network: network2,
			},
			args: args{
				testProject,
				[]bool{},
			},
			want: want{
				err:     errors.OK,
				network: network2,
			},
		},
		{
			name: "Update existing network (update)",
			fields: fields{
				network: network3,
			},
			args: args{
				testProject,
				[]bool{true},
			},
			want: want{
				err:     errors.OK,
				network: network3,
			},
		},
		{
			name: "Update existing network (no update)",
			fields: fields{
				network: network4,
			},
			args: args{
				testProject,
				[]bool{},
			},
			want: want{
				err:     errors.Conflict,
				network: network3,
			},
		},
		{
			name: "Update non existing network",
			fields: fields{
				network: network5,
			},
			args: args{
				testProject,
				[]bool{true},
			},
			want: want{
				err:     errors.NotFound,
				network: Network{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network := tt.fields.network
			if got := network.Insert(tt.args.project, tt.args.update...); got.Code != tt.want.err.Code {
				t.Errorf("Insert() = %v, want %v", got.Code, tt.want.err.Code)
			}
			id := tt.fields.network.Identifier
			if !reflect.DeepEqual(testProject.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.ID], tt.want.network) {
				t.Errorf("Insert() = %v, want %v", network, tt.want.network)
			}
		})
	}
}

func TestNetwork_ToDomain(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Network
		KubernetesResources kubernetes.ResourceList
		Subnetworks         SubnetworkCollection
		Firewalls           FirewallCollection
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
			network := &Network{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Subnetworks:         tt.fields.Subnetworks,
				Firewalls:           tt.fields.Firewalls,
			}
			got, got1 := network.ToDomain()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDomain() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ToDomain() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNetwork_WithMetadata(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Network
		KubernetesResources kubernetes.ResourceList
		Subnetworks         SubnetworkCollection
		Firewalls           FirewallCollection
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
			network := &Network{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Subnetworks:         tt.fields.Subnetworks,
				Firewalls:           tt.fields.Firewalls,
			}
			network.WithMetadata(tt.args.request)
		})
	}
}

func TestNewNetwork(t *testing.T) {
	type args struct {
		id identifier.Network
	}
	tests := []struct {
		name string
		args args
		want Network
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNetwork(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}
