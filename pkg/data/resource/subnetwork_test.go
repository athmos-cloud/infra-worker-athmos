package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
	"testing"
)

func TestSubnetwork_FromMap(t *testing.T) {
	type fields struct {
		Subnetwork Subnetwork
	}
	type args struct {
		data map[string]interface{}
	}
	type want struct {
		err        errors.Error
		subnetwork Subnetwork
	}
	subnet := NewSubnetwork(identifier.Subnetwork{ID: "test", ProviderID: "test", NetworkID: "test"})
	expectedSubnet := subnet
	expectedSubnet.Region = "eu-west-1"
	expectedSubnet.IPCIDRRange = "10.0.0.0/8"

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "FromMap with valid data",
			fields: fields{
				Subnetwork: subnet,
			},
			args: args{
				data: map[string]interface{}{
					"region":      "eu-west-1",
					"ipCidrRange": "10.0.0.0/8",
				},
			},
			want: want{
				err:        errors.OK,
				subnetwork: expectedSubnet,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curVpc := tt.fields.Subnetwork
			if got := curVpc.FromMap(tt.args.data); !reflect.DeepEqual(got.Code, tt.want.err.Code) {
				t.Errorf("FromMap() = %v, want %v", got, tt.want)
			}
			if !curVpc.Equals(tt.want.subnetwork) {
				t.Errorf("FromMap() = %v, want %v", curVpc, tt.want.subnetwork)
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
			subnetworkGot := testProject.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.ID]
			if !subnetworkGot.Equals(tt.want.subnetwork) {
				t.Errorf("Insert() = %v, want %v", subnet, tt.want.subnetwork)
			}
		})
	}
}
