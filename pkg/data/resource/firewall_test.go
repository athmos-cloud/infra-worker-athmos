package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
	"testing"
)

func TestFirewall_FromMap(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Firewall
		KubernetesResources kubernetes.ResourceList
		Network             string
		Allow               RuleList
		Deny                RuleList
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
			firewall := &Firewall{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Network:             tt.fields.Network,
				Allow:               tt.fields.Allow,
				Deny:                tt.fields.Deny,
			}
			if got := firewall.FromMap(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirewall_GetMetadata(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Firewall
		KubernetesResources kubernetes.ResourceList
		Network             string
		Allow               RuleList
		Deny                RuleList
	}
	tests := []struct {
		name   string
		fields fields
		want   metadata.Metadata
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firewall := &Firewall{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Network:             tt.fields.Network,
				Allow:               tt.fields.Allow,
				Deny:                tt.fields.Deny,
			}
			if got := firewall.GetMetadata(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirewall_GetPluginReference(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Firewall
		KubernetesResources kubernetes.ResourceList
		Network             string
		Allow               RuleList
		Deny                RuleList
	}
	type args struct {
		request resource.GetPluginReferenceRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   resource.GetPluginReferenceResponse
		want1  errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firewall := &Firewall{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Network:             tt.fields.Network,
				Allow:               tt.fields.Allow,
				Deny:                tt.fields.Deny,
			}
			got, got1 := firewall.GetPluginReference(tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPluginReference() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetPluginReference() got1 = %v, want %v", got1, tt.want1)
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
			if !reflect.DeepEqual(testProject.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Firewalls[id.ID], tt.want.firewall) {
				t.Errorf("Insert() = %v, want %v", firewall, tt.want.firewall)
			}
		})
	}
}

func TestFirewall_ToDomain(t *testing.T) {
	type fields struct {
		Firewall Firewall
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
			firewall := tt.fields.Firewall
			got, got1 := firewall.ToDomain()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDomain() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ToDomain() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFirewall_WithMetadata(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Firewall
		KubernetesResources kubernetes.ResourceList
		Network             string
		Allow               RuleList
		Deny                RuleList
	}
	type args struct {
		request metadata.CreateMetadataRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firewall := &Firewall{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Network:             tt.fields.Network,
				Allow:               tt.fields.Allow,
				Deny:                tt.fields.Deny,
			}
			firewall.WithMetadata(tt.args.request)
		})
	}
}

func TestNewFirewall(t *testing.T) {
	type args struct {
		id identifier.Firewall
	}
	tests := []struct {
		name string
		args args
		want Firewall
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFirewall(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFirewall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuleList_FromMap(t *testing.T) {
	type args struct {
		data []interface{}
	}
	tests := []struct {
		name  string
		rules RuleList
		args  args
		want  errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rules.FromMap(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
