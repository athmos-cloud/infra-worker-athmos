package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
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
	subnet := NewSubnetwork(NewResourcePayload{
		Name: "test",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: "test",
			VPCID:      "test",
			NetworkID:  "test",
		}),
		Provider: types.GCP,
	})
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
		update  bool
	}
	type want struct {
		err        errors.Error
		subnetwork Subnetwork
	}
	providerID := "test"
	vpcID := "test"
	networkID := "test"

	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(NewResourcePayload{
		Name:             providerID,
		ParentIdentifier: identifier.Empty{},
		Provider:         types.GCP,
	})
	testVPC := NewVPC(NewResourcePayload{
		Name: vpcID,
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
		}),
		Provider: types.GCP,
	})
	testNetwork := NewNetwork(NewResourcePayload{
		Name: networkID,
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
		}),
		Provider: types.GCP,
	})
	_ = identifier.Build(identifier.IdPayload{
		ProviderID: providerID,
		VPCID:      vpcID,
		NetworkID:  networkID,
	})
	subnetwork1 := NewSubnetwork(NewResourcePayload{
		Name: "test-1",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
			NetworkID:  networkID,
		}),
		Provider: types.GCP,
	})
	subnetwork1.Identifier.SubnetworkID = "test-1"
	subnetwork2 := NewSubnetwork(NewResourcePayload{
		Name: "test-1",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
			NetworkID:  networkID,
		}),
		Provider: types.GCP,
	})
	subnetwork2.Identifier.SubnetworkID = "test-2"
	subnetwork3 := subnetwork1
	subnetwork3.Metadata.Tags = map[string]string{"test": "test"}
	subnetwork4 := subnetwork3
	subnetwork4.Metadata.Tags = map[string]string{"hello": "world"}
	subnetwork5 := NewSubnetwork(NewResourcePayload{
		Name: "test-5",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
			NetworkID:  networkID,
		}),
		Provider: types.GCP,
	})
	subnetwork5.Identifier.SubnetworkID = "test-5"

	testNetwork.Subnetworks[subnetwork1.Identifier.SubnetworkID] = subnetwork1
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
				*testProject,
				false,
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
				*testProject,
				true,
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
				*testProject,
				false,
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
				*testProject,
				true,
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
			if tt.args.update {
				tt.args.project.Update(&subnet)
			} else {
				tt.args.project.Insert(&subnet)
			}
			id := tt.fields.subnetwork.Identifier
			subnetworkGot := tt.args.project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.SubnetworkID]
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
	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(NewResourcePayload{
		Name:             providerID,
		ParentIdentifier: identifier.Empty{},
		Provider:         types.GCP,
	})
	testVPC := NewVPC(NewResourcePayload{
		Name: vpcID,
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
		}),
		Provider: types.GCP,
	})
	testNetwork := NewNetwork(NewResourcePayload{
		Name: networkID,
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
		}),
		Provider: types.GCP,
	})
	subnetwork1 := NewSubnetwork(NewResourcePayload{
		Name: "test-1",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
			NetworkID:  networkID,
		}),
		Provider: types.GCP,
	})
	subnetwork2 := NewSubnetwork(NewResourcePayload{
		Name: "test-1",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
			NetworkID:  networkID,
		}),
		Provider: types.GCP,
	})

	testNetwork.Subnetworks[subnetwork1.Identifier.SubnetworkID] = subnetwork1
	testVPC.Networks[networkID] = testNetwork
	testProvider.VPCs[vpcID] = testVPC
	testProject.Resources[providerID] = testProvider
	testNetwork.Subnetworks[subnetwork1.Identifier.SubnetworkID] = subnetwork1
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
				project: *testProject,
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
				project: *testProject,
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
			tt.args.project.Delete(&subnet)
		})
	}
}
