package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	types2 "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"testing"
)

func TestVM_FromMap(t *testing.T) {
	type fields struct {
		VM VM
	}
	type args struct {
		data map[string]interface{}
	}
	type want struct {
		err errors.Error
		vm  VM
	}
	vmID := identifier.Build(identifier.IdPayload{
		ProviderID: "test",
		NetworkID:  "test",
		SubnetID:   "test",
	})
	vm := NewVM(NewResourcePayload{
		Name:             "test",
		ParentIdentifier: vmID,
		Provider:         types2.GCP,
	})
	vm.Identifier.VMID = "test"
	expectedVM1 := vm
	expectedVM1.Zone = "europe-west1-a"
	expectedVM1.MachineType = "f1-micro"
	expectedVM1.Disks = []Disk{
		{
			Type:       "SSD",
			Mode:       types2.ReadOnly,
			SizeGib:    10,
			AutoDelete: true,
		},
	}
	expectedVM1.Auths = []VMAuth{
		{
			Username:     "admin",
			SSHPublicKey: "cfrezverververvre",
		},
	}
	expectedVM2 := vm

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "FromMap with valid data",
			fields: fields{
				VM: vm,
			},
			args: args{
				data: map[string]interface{}{
					"vpc":         "vpc-test",
					"zone":        "europe-west1-a",
					"machineType": "f1-micro",
					"disk": []map[string]interface{}{
						{
							"type":       "SSD",
							"diskMode":   "READ_ONLY",
							"sizeGib":    10,
							"autoDelete": true,
						},
					},
					"auths": []map[string]interface{}{
						{
							"username":     "admin",
							"sshPublicKey": "cfrezverververvre",
						},
					},
				},
			},
			want: want{
				vm: expectedVM1,
			},
		}, {
			name: "FromMap with invalid data",
			fields: fields{
				VM: vm,
			},
			args: args{
				data: map[string]interface{}{
					"disk": map[string]interface{}{
						"type":       "SSD",
						"diskMode":   "wakanda",
						"sizeGib":    10,
						"autoDelete": true,
					},
				},
			},
			want: want{
				err: errors.InvalidArgument,
				vm:  expectedVM2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curVM := tt.fields.VM
			defer func() {
				r := recover()
				if r != nil {
					err := r.(errors.Error)
					if err.Code != tt.want.err.Code {
						t.Errorf("FromMap()  %v, want %v", err.Code, tt.want.err.Code)
					}
				}
			}()
			curVM.FromMap(tt.args.data)
			if !curVM.Equals(tt.want.vm) {
				t.Errorf("FromMap() = %v, want %v", curVM.Auths, tt.want.vm.Auths)
			}
		})
	}
}

func TestVM_Insert(t *testing.T) {
	type fields struct {
		vm VM
	}
	type args struct {
		project Project
		update  bool
	}
	type want struct {
		err errors.Error
		vm  VM
	}
	providerID := "test"
	vpcID := "test"
	networkID := "test"
	subnetID := "test"

	vm1 := NewVM(NewResourcePayload{
		Name: "test-1",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			NetworkID:  networkID,
			SubnetID:   subnetID,
		}),
		Provider: types2.GCP,
	})
	vm1.Identifier.VMID = "test-1"
	vm2 := NewVM(NewResourcePayload{
		Name: "test-2",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			NetworkID:  networkID,
			SubnetID:   subnetID,
		}),
		Provider: types2.GCP,
	})
	vm2.Identifier.VMID = "test-2"
	vm3 := vm1
	vm3.Metadata.Tags = map[string]string{"test": "test"}
	vm4 := vm3
	vm4.Metadata.Tags = map[string]string{"hello": "world"}
	vm5 := NewVM(NewResourcePayload{
		Name: "test-5",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			NetworkID:  networkID,
			SubnetID:   subnetID,
		}),
		Provider: types2.GCP,
	})

	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(NewResourcePayload{
		Name:             providerID,
		ParentIdentifier: identifier.Empty{},
		Provider:         types2.GCP,
	})
	testVPC := NewVPC(NewResourcePayload{
		Name: vpcID,
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
		}),
		Provider: types2.GCP,
	})
	testNetwork := NewNetwork(NewResourcePayload{
		Name: networkID,
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
		}),
		Provider: types2.GCP,
	})
	testSubnet := NewSubnetwork(NewResourcePayload{
		Name: networkID,
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
			NetworkID:  networkID,
		}),
		Provider: types2.GCP,
	})
	testSubnet.VMs[vm1.Identifier.VMID] = vm1
	testNetwork.Subnetworks[subnetID] = testSubnet
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
			name: "Insert non existing vm (creation)",
			fields: fields{
				vm: vm2,
			},
			args: args{
				*testProject,
				false,
			},
			want: want{
				vm: vm2,
			},
		},
		{
			name: "Update existing vm (update)",
			fields: fields{
				vm: vm3,
			},
			args: args{
				*testProject,
				true,
			},
			want: want{
				vm: vm3,
			},
		},
		{
			name: "Update existing vm (no update)",
			fields: fields{
				vm: vm4,
			},
			args: args{
				*testProject,
				false,
			},
			want: want{
				err: errors.Conflict,
				vm:  vm3,
			},
		},
		{
			name: "Update non existing vm",
			fields: fields{
				vm: vm5,
			},
			args: args{
				*testProject,
				true,
			},
			want: want{
				err: errors.NotFound,
				vm:  VM{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := tt.fields.vm
			defer func() {
				if r := recover(); r != nil {
					err := r.(errors.Error)
					if err.Code != tt.want.err.Code {
						t.Errorf("Insert()  %v, want %v", err.Code, tt.want.err.Code)
					}
				}
			}()
			if tt.args.update {
				tt.args.project.Update(&vm)
			} else {
				tt.args.project.Insert(&vm)
			}
			id := tt.fields.vm.Identifier
			gotVM := tt.args.project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.SubnetID].VMs[id.VMID]
			if !gotVM.Equals(tt.want.vm) {
				t.Errorf("Insert() = %v, want %v", vm, tt.want.vm)
			}
		})
	}
}

func TestVM_Remove(t *testing.T) {
	type fields struct {
		vm VM
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
	subnetID := "test"

	vm1 := NewVM(NewResourcePayload{
		Name: "test-1",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			NetworkID:  networkID,
			SubnetID:   subnetID,
		}),
		Provider: types2.GCP,
	})
	vm2 := NewVM(NewResourcePayload{
		Name: "test-2",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			NetworkID:  networkID,
			SubnetID:   subnetID,
		}),
		Provider: types2.GCP,
	})
	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(NewResourcePayload{
		Name:             providerID,
		ParentIdentifier: identifier.Empty{},
		Provider:         types2.GCP,
	})
	testVPC := NewVPC(NewResourcePayload{
		Name: vpcID,
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
		}),
		Provider: types2.GCP,
	})
	testNetwork := NewNetwork(NewResourcePayload{
		Name: networkID,
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
		}),
		Provider: types2.GCP,
	})
	testSubnet := NewSubnetwork(NewResourcePayload{
		Name: networkID,
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: providerID,
			VPCID:      vpcID,
			NetworkID:  networkID,
		}),
		Provider: types2.GCP,
	})
	testSubnet.VMs[vm1.Identifier.VMID] = vm1
	testNetwork.Subnetworks[subnetID] = testSubnet
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
			name: "Remove existing vm",
			fields: fields{
				vm: vm1,
			},
			args: args{
				*testProject,
			},
			want: want{},
		},
		{
			name: "Remove non-existing vm",
			fields: fields{
				vm: vm2,
			},
			args: args{
				*testProject,
			},
			want: want{
				err: errors.NotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := tt.fields.vm
			defer func() {
				if r := recover(); r != nil {
					err := r.(errors.Error)
					if err.Code != tt.want.err.Code {
						t.Errorf("FromMap()  %v, want %v", err.Code, tt.want.err.Code)
					}
				}
			}()
			tt.args.project.Delete(&vm)
		})
	}
}
