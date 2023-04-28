package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"reflect"
	"testing"
)

func TestFirewall_FromMap(t *testing.T) {
	type fields struct {
		Firewall Firewall
	}
	type args struct {
		data map[string]interface{}
	}
	type want struct {
		err      errors.Error
		firewall Firewall
	}
	firewall := NewFirewall(identifier.Firewall{ID: "test", ProviderID: "test", NetworkID: "test"})
	expectedFirewall := firewall
	expectedFirewall.Allow = RuleList{
		{
			Protocol: "tcp",
			Ports:    []int{22},
		},
		{
			Protocol: "udp",
			Ports:    []int{80, 8080},
		},
	}
	expectedFirewall.Deny = RuleList{
		{
			Protocol: "tcp",
			Ports:    []int{222},
		},
		{
			Protocol: "udp",
			Ports:    []int{81},
		},
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "FromMap with valid data",
			fields: fields{
				Firewall: firewall,
			},
			args: args{
				data: map[string]interface{}{
					"allow": []map[string]interface{}{
						{
							"protocol": "tcp",
							"ports":    []interface{}{22},
						},
						{
							"protocol": "udp",
							"ports":    []interface{}{80, 8080},
						},
					},
					"deny": []map[string]interface{}{
						{
							"protocol": "tcp",
							"ports":    []interface{}{222},
						},
						{
							"protocol": "udp",
							"ports":    []interface{}{81},
						},
					},
				},
			},
			want: want{
				err:      errors.OK,
				firewall: expectedFirewall,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curFirewall := tt.fields.Firewall
			if got := curFirewall.FromMap(tt.args.data); !reflect.DeepEqual(got.Code, tt.want.err.Code) {
				logger.Info.Println(got)
				t.Errorf("FromMap() = %v, want %v", got.Code, tt.want.err.Code)
			}
			if !curFirewall.Equals(tt.want.firewall) {
				t.Errorf("FromMap() = %v, want %v", curFirewall, tt.want.firewall)
			}
		})
	}
}

func TestFirewall_Insert(t *testing.T) {
	type fields struct {
		firewall Firewall
	}
	type args struct {
		project Project
		update  []bool
	}
	type want struct {
		err      errors.Error
		firewall Firewall
	}
	providerID := "test"
	vpcID := "test"
	networkID := "test"

	firewall1 := NewFirewall(identifier.Firewall{ID: "test-1", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID})
	firewall2 := NewFirewall(identifier.Firewall{ID: "test-2", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID})
	firewall3 := firewall1
	firewall3.Metadata.Tags = map[string]string{"test": "test"}
	firewall4 := firewall3
	firewall4.Metadata.Tags = map[string]string{"hello": "world"}
	firewall5 := NewFirewall(identifier.Firewall{ID: "test-5", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID})

	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(identifier.Provider{ID: providerID})
	testVPC := NewVPC(identifier.VPC{ID: vpcID, ProviderID: providerID})
	testNetwork := NewNetwork(identifier.Network{ID: networkID, ProviderID: providerID, VPCID: vpcID})
	testNetwork.Firewalls[firewall1.Identifier.ID] = firewall1
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
			name: "Insert non existing firewall (creation)",
			fields: fields{
				firewall: firewall2,
			},
			args: args{
				testProject,
				[]bool{},
			},
			want: want{
				err:      errors.OK,
				firewall: firewall2,
			},
		},
		{
			name: "Update existing firewall (update)",
			fields: fields{
				firewall: firewall3,
			},
			args: args{
				testProject,
				[]bool{true},
			},
			want: want{
				err:      errors.OK,
				firewall: firewall3,
			},
		},
		{
			name: "Update existing firewall (no update)",
			fields: fields{
				firewall: firewall4,
			},
			args: args{
				testProject,
				[]bool{},
			},
			want: want{
				err:      errors.Conflict,
				firewall: firewall3,
			},
		},
		{
			name: "Update non existing firewall",
			fields: fields{
				firewall: firewall5,
			},
			args: args{
				testProject,
				[]bool{true},
			},
			want: want{
				err:      errors.NotFound,
				firewall: Firewall{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firewall := tt.fields.firewall
			if got := firewall.Insert(tt.args.project, tt.args.update...); !reflect.DeepEqual(got.Code, tt.want.err.Code) {
				t.Errorf("Insert() = %v, want %v", got.Code, tt.want.err.Code)
			}
			id := tt.fields.firewall.Identifier
			firewallGot := testProject.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Firewalls[id.ID]
			if !firewallGot.Equals(tt.want.firewall) {
				t.Errorf("Insert() = %v, want %v", firewall, tt.want.firewall)
			}
		})
	}
}

func TestFirewall_Remove(t *testing.T) {
	type fields struct {
		Firewall Firewall
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

	firewall1 := NewFirewall(identifier.Firewall{ID: "test-1", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID})
	firewall2 := NewFirewall(identifier.Firewall{ID: "test-2", ProviderID: providerID, VPCID: vpcID, NetworkID: networkID})

	testProject := NewProject("test", "owner_test")
	testProvider := NewProvider(identifier.Provider{ID: providerID})
	testVPC := NewVPC(identifier.VPC{ID: vpcID, ProviderID: providerID})
	testNetwork := NewNetwork(identifier.Network{ID: networkID, ProviderID: providerID, VPCID: vpcID})
	testNetwork.Firewalls[firewall1.Identifier.ID] = firewall1
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
			name: "Remove existing firewall",
			fields: fields{
				Firewall: firewall1,
			},
			args: args{
				project: testProject,
			},
			want: want{
				err: errors.NoContent,
			},
		},
		{
			name: "Remove non-existing firewall",
			fields: fields{
				Firewall: firewall2,
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
			firewall := tt.fields.Firewall
			if got := firewall.Remove(tt.args.project); !reflect.DeepEqual(got.Code, tt.want.err.Code) {
				t.Errorf("Remove() = %v, want %v", got.Code, tt.want.err.Code)
			}
		})
	}
}
