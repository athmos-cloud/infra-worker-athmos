package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/types"
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
	vm := NewVM(identifier.VM{ID: "test", SubnetID: "test", NetworkID: "test", VPCID: "test", ProviderID: "test"}, common.Azure)
	expectedVM1 := vm
	expectedVM1.VPC = "vpc-test"
	expectedVM1.Zone = "europe-west1-a"
	expectedVM1.MachineType = "f1-micro"
	expectedVM1.Disk = Disk{
		Type:       "SSD",
		Mode:       types.ReadOnly,
		SizeGib:    10,
		AutoDelete: true,
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
					"disks": []map[string]interface{}{
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
					"disks": []map[string]interface{}{
						{
							"type":       "SSD",
							"diskMode":   "wakanda",
							"sizeGib":    10,
							"autoDelete": true,
						},
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
				if r := recover(); r != nil {
					err := r.(errors.Error)
					if err.Code != tt.want.err.Code {
						t.Errorf("FromMap()  %v, want %v", err.Code, tt.want.err.Code)
					}
				}
			}()
			curVM.FromMap(tt.args.data)
			if !curVM.Equals(tt.want.vm) {
				t.Errorf("FromMap() = %v, want %v", curVM, tt.want.vm)
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
		update  []bool
	}
	type want struct {
		err errors.Error
		vm  VM
	}
	providerID := "test"
	vpcID := "test"
	networkID := "test"
	subnetID := "test"

	vm1 := NewVM(identifier.VM{ID: "test-1", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID, SubnetID: subnetID}, common.Azure)
	vm2 := NewVM(identifier.VM{ID: "test-2", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID, SubnetID: subnetID}, common.Azure)
	vm3 := vm1
	vm3.Metadata.Tags = map[string]string{"test": "test"}
	vm4 := vm3
	vm4.Metadata.Tags = map[string]string{"hello": "world"}
	vm5 := NewVM(identifier.VM{ID: "test-5", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID, SubnetID: subnetID}, common.Azure)

	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(identifier.Provider{ID: providerID}, common.Azure)
	testVPC := NewVPC(identifier.VPC{ID: vpcID, ProviderID: providerID}, common.Azure)
	testNetwork := NewNetwork(identifier.Network{ID: networkID, ProviderID: providerID, VPCID: vpcID}, common.Azure)
	testSubnet := NewSubnetwork(identifier.Subnetwork{ID: subnetID, ProviderID: providerID, VPCID: vpcID, NetworkID: networkID}, common.Azure)
	testSubnet.VMs[vm1.Identifier.ID] = vm1
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
				testProject,
				[]bool{},
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
				testProject,
				[]bool{true},
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
				testProject,
				[]bool{},
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
				testProject,
				[]bool{true},
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
						t.Errorf("FromMap()  %v, want %v", err.Code, tt.want.err.Code)
					}
				}
			}()
			vm.Insert(tt.args.project, tt.args.update...)
			id := tt.fields.vm.Identifier
			gotVM := testProject.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.SubnetID].VMs[id.ID]
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

	vm1 := NewVM(identifier.VM{ID: "test-1", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID, SubnetID: subnetID}, common.Azure)
	vm2 := NewVM(identifier.VM{ID: "test-2", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID, SubnetID: subnetID}, common.Azure)

	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(identifier.Provider{ID: providerID}, common.Azure)
	testVPC := NewVPC(identifier.VPC{ID: vpcID, ProviderID: providerID}, common.Azure)
	testNetwork := NewNetwork(identifier.Network{ID: networkID, ProviderID: providerID, VPCID: vpcID}, common.Azure)
	testSubnet := NewSubnetwork(identifier.Subnetwork{ID: subnetID, ProviderID: providerID, VPCID: vpcID, NetworkID: networkID}, common.Azure)
	testSubnet.VMs[vm1.Identifier.ID] = vm1
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
				testProject,
			},
			want: want{},
		},
		{
			name: "Remove non-existing vm",
			fields: fields{
				vm: vm2,
			},
			args: args{
				testProject,
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
			vm.Remove(tt.args.project)
		})
	}
}
