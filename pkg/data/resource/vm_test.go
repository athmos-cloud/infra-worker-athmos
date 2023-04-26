package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
	"testing"
)

func TestVM_FromMap(t *testing.T) {
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
	type args struct {
		data map[string]interface{}
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
			if got := vm.FromMap(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromMap() = %v, want %v", got, tt.want)
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
			if !reflect.DeepEqual(testProject.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.SubnetID].VMs[id.ID], tt.want.vm) {
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
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ToDomain() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
