package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"testing"
)

func TestNetwork_Insert(t *testing.T) {
	type fields struct {
		network Network
	}
	type args struct {
		project Project
		update  bool
	}
	type want struct {
		err     errors.Error
		network Network
	}

	providerID := "test"
	vpcID := "test"

	network1 := NewNetwork(NewResourcePayload{
		Name: "test-1",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
		}),
		Provider: types.GCP,
	})
	network1.Identifier.NetworkID = "test"
	network2 := NewNetwork(NewResourcePayload{
		Name: "test-2",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
		}),
		Provider: types.GCP,
	})
	network2.Identifier.NetworkID = "test-2"
	network3 := network1
	network3.Metadata.Tags = map[string]string{"test": "test"}
	network4 := network3
	network4.Metadata.Tags = map[string]string{"hello": "world"}
	network5 := NewNetwork(NewResourcePayload{
		Name: "test-5",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
		}),
		Provider: types.GCP,
	})

	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(NewResourcePayload{
		Name:             providerID,
		ParentIdentifier: identifier.Empty{},
		Provider:         types.GCP,
	})
	testProvider.Identifier.ProviderID = providerID
	testVPC := NewVPC(NewResourcePayload{
		Name: vpcID,
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
		}),
		Provider: types.GCP,
	})
	testVPC.Identifier.VPCID = vpcID
	testVPC.Networks[network1.Identifier.NetworkID] = network1
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
				*testProject,
				false,
			},
			want: want{
				network: network2,
			},
		},
		{
			name: "Update existing network (update)",
			fields: fields{
				network: network3,
			},
			args: args{
				*testProject,
				true,
			},
			want: want{
				network: network3,
			},
		},
		{
			name: "Update existing network (no update)",
			fields: fields{
				network: network4,
			},
			args: args{
				*testProject,
				true,
			},
			want: want{
				err:     errors.Conflict,
				network: network4,
			},
		},
		{
			name: "Update non existing network",
			fields: fields{
				network: network5,
			},
			args: args{
				*testProject,
				true,
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
			defer func() {
				if r := recover(); r != nil {
					err := r.(errors.Error)
					if err.Code != tt.want.err.Code {
						t.Errorf("FromMap()  %v, want %v", err.Code, tt.want.err.Code)
					}
				}
			}()

			if tt.args.update {
				tt.args.project.Update(&network)
			} else {
				tt.args.project.Insert(&network)
			}
			id := tt.fields.network.Identifier
			networkGot := tt.args.project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID]
			if !networkGot.Equals(tt.want.network) {
				t.Errorf("Insert() = %v, want %v", network, tt.want.network)
			}
		})
	}
}

func TestNetwork_Remove(t *testing.T) {
	type fields struct {
		Network Network
	}
	type args struct {
		project Project
	}
	type want struct {
		err errors.Error
	}

	providerID := "test"
	vpcID := "test"

	network1 := NewNetwork(NewResourcePayload{
		Name: "test-1",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
		}),
		Provider: types.GCP,
	})
	network2 := NewNetwork(NewResourcePayload{
		Name: "test-2",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
		}),
		Provider: types.GCP,
	})

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
	testVPC.Networks[network1.Identifier.NetworkID] = network1
	testProvider.VPCs[vpcID] = testVPC
	testProject.Resources[providerID] = testProvider

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Remove existing network",
			fields: fields{
				Network: network1,
			},
			args: args{
				project: *testProject,
			},
			want: want{
				err: errors.NoContent,
			},
		},
		{
			name: "Remove non-existing network",
			fields: fields{
				Network: network2,
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
			network := tt.fields.Network
			defer func() {
				if r := recover(); r != nil {
					err := r.(errors.Error)
					if err.Code != tt.want.err.Code {
						t.Errorf("FromMap()  %v, want %v", err.Code, tt.want.err.Code)
					}
				}
			}()
			tt.args.project.Delete(&network)
		})
	}
}
