package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"reflect"
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
	vm := NewVM(identifier.VM{ID: "test", SubnetID: "test", NetworkID: "test", VPCID: "test", ProviderID: "test"})
	expectedVM1 := vm
	expectedVM1.VPC = "vpc-test"
	expectedVM1.Zone = "europe-west1-a"
	expectedVM1.MachineType = "f1-micro"
	expectedVM1.Disks = []Disk{
		{
			Type:       "SSD",
			Mode:       types.ReadOnly,
			SizeGib:    10,
			AutoDelete: true,
		},
		{
			Type:       "SSD",
			Mode:       types.ReadWrite,
			SizeGib:    100,
			AutoDelete: false,
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
					"disks": []map[string]interface{}{
						{
							"type":       "SSD",
							"diskMode":   "READ_ONLY",
							"sizeGib":    10,
							"autoDelete": true,
						},
						{
							"type":     "SSD",
							"diskMode": "READ_WRITE",
							"sizeGib":  100,
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
				err: errors.OK,
				vm:  expectedVM1,
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
			got := curVM.FromMap(tt.args.data)
			if got.Code != tt.want.err.Code {
				logger.Info.Println(got)
				t.Errorf("FromMap() = %v, want %v", got.Code, tt.want.err.Code)
				return
			}
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

	vm1 := NewVM(identifier.VM{ID: "test-1", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID, SubnetID: subnetID})
	vm2 := NewVM(identifier.VM{ID: "test-2", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID, SubnetID: subnetID})
	vm3 := vm1
	vm3.Metadata.Tags = map[string]string{"test": "test"}
	vm4 := vm3
	vm4.Metadata.Tags = map[string]string{"hello": "world"}
	vm5 := NewVM(identifier.VM{ID: "test-5", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID, SubnetID: subnetID})

	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(identifier.Provider{ID: providerID})
	testVPC := NewVPC(identifier.VPC{ID: vpcID, ProviderID: providerID})
	testNetwork := NewNetwork(identifier.Network{ID: networkID, ProviderID: providerID, VPCID: vpcID})
	testSubnet := NewSubnetwork(identifier.Subnetwork{ID: subnetID, ProviderID: providerID, VPCID: vpcID, NetworkID: networkID})
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
				err: errors.OK,
				vm:  vm2,
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
				err: errors.OK,
				vm:  vm3,
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
			if got := vm.Insert(tt.args.project, tt.args.update...); !reflect.DeepEqual(got.Code, tt.want.err.Code) {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
			id := tt.fields.vm.Identifier
			gotVM := testProject.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.SubnetID].VMs[id.ID]
			if !gotVM.Equals(tt.want.vm) {
				t.Errorf("Insert() = %v, want %v", vm, tt.want.vm)
			}
		})
	}
}

func TestVM_ToDomain(t *testing.T) {
	type fields struct {
		Metadata    metadata.Metadata
		Identifier  identifier.VM
		VPC         string
		Network     string
		Subnetwork  string
		Zone        string
		MachineType string
		Auths       []VMAuth
		Disks       []Disk
		OS          OS
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
			vm := &VM{
				Metadata:    tt.fields.Metadata,
				Identifier:  tt.fields.Identifier,
				VPC:         tt.fields.VPC,
				Network:     tt.fields.Network,
				Subnetwork:  tt.fields.Subnetwork,
				Zone:        tt.fields.Zone,
				MachineType: tt.fields.MachineType,
				Auths:       tt.fields.Auths,
				Disks:       tt.fields.Disks,
				OS:          tt.fields.OS,
			}
			got, got1 := vm.ToDomain()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDomain() got = %v, want %v", got, tt.want)
			}
			if !got1.Equals(tt.want1) {
				t.Errorf("ToDomain() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
