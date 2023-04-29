package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
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
	subnet := NewSubnetwork(identifier.Subnetwork{ID: "test", ProviderID: "test", NetworkID: "test"}, common.GCP)
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
				subnetwork: expectedSubnet,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curVpc := tt.fields.Subnetwork
			defer func() {
				if r := recover(); r != nil {
					err := r.(errors.Error)
					if err.Code != tt.want.err.Code {
						t.Errorf("FromMap()  %v, want %v", err.Code, tt.want.err.Code)
					}
				}
			}()
			curVpc.FromMap(tt.args.data)
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

	subnetwork1 := NewSubnetwork(identifier.Subnetwork{ID: "test-1", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID}, common.GCP)
	subnetwork2 := NewSubnetwork(identifier.Subnetwork{ID: "test-2", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID}, common.GCP)
	subnetwork3 := subnetwork1
	subnetwork3.Metadata.Tags = map[string]string{"test": "test"}
	subnetwork4 := subnetwork3
	subnetwork4.Metadata.Tags = map[string]string{"hello": "world"}
	subnetwork5 := NewSubnetwork(identifier.Subnetwork{ID: "test-5", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID}, common.GCP)

	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(identifier.Provider{ID: providerID}, common.GCP)
	testVPC := NewVPC(identifier.VPC{ID: vpcID, ProviderID: providerID}, common.GCP)
	testNetwork := NewNetwork(identifier.Network{ID: networkID, ProviderID: providerID, VPCID: vpcID}, common.GCP)
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
			defer func() {
				if r := recover(); r != nil {
					err := r.(errors.Error)
					if err.Code != tt.want.err.Code {
						t.Errorf("FromMap()  %v, want %v", err.Code, tt.want.err.Code)
					}
				}
			}()
			subnet.Insert(tt.args.project, tt.args.update...)
			id := tt.fields.subnetwork.Identifier
			subnetworkGot := testProject.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.ID]
			if !subnetworkGot.Equals(tt.want.subnetwork) {
				t.Errorf("Insert() = %v, want %v", subnet, tt.want.subnetwork)
			}
		})
	}
}

func TestSubnetwork_Remove(t *testing.T) {
	type fields struct {
		Subnetwork Subnetwork
	}
	type args struct {
		project Project
	}
	type want struct {
		err errors.Error
	}

	providerID := "test"
	vpcID := "test"
	networkID := "test"

	subnetwork1 := NewSubnetwork(identifier.Subnetwork{ID: "test-1", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID}, common.GCP)
	subnetwork2 := NewSubnetwork(identifier.Subnetwork{ID: "test-2", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID}, common.GCP)

	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(identifier.Provider{ID: providerID}, common.GCP)
	testVPC := NewVPC(identifier.VPC{ID: vpcID, ProviderID: providerID}, common.GCP)
	testNetwork := NewNetwork(identifier.Network{ID: networkID, ProviderID: providerID, VPCID: vpcID}, common.GCP)
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
			name: "Remove existing subnetwork",
			fields: fields{
				Subnetwork: subnetwork1,
			},
			args: args{
				project: testProject,
			},
			want: want{
				err: errors.NoContent,
			},
		},
		{
			name: "Remove non-existing subnetwork",
			fields: fields{
				Subnetwork: subnetwork2,
			},
			args: args{
				project: testProject,
			},
			want: want{
				err: errors.NotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subnet := tt.fields.Subnetwork
			defer func() {
				if r := recover(); r != nil {
					err := r.(errors.Error)
					if err.Code != tt.want.err.Code {
						t.Errorf("FromMap()  %v, want %v", err.Code, tt.want.err.Code)
					}
				}
			}()
			subnet.Remove(tt.args.project)
		})
	}
}
