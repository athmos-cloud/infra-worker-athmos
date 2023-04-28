package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"testing"
)

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
			networkGot := testProject.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.ID]
			if !networkGot.Equals(tt.want.network) {
				t.Errorf("Insert() = %v, want %v", network, tt.want.network)
			}
		})
	}
}
