package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	auth "github.com/athmos-cloud/infra-worker-athmos/pkg/data/auth"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
	"testing"
)

func TestProvider_FromMap(t *testing.T) {
	type fields struct {
		Provider Provider
	}
	type args struct {
		m map[string]interface{}
	}
	type want struct {
		err      errors.Error
		provider Provider
	}
	provider := NewProvider(identifier.Provider{ID: "test"})
	expectedProvider1 := provider
	expectedProvider1.Auth = auth.Auth{
		AuthType: auth.AuthTypeSecret,
		SecretAuth: auth.SecretAuth{
			SecretName: "key-secret",
			SecretKey:  "key.json",
			Namespace:  "default",
		},
	}
	expectedProvider1.VPC = "vpc-test"
	expectedProvider2 := provider
	expectedProvider2.VPC = "vpc-test"

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "FromMap with valid data",
			fields: fields{
				Provider: provider,
			},
			args: args{
				m: map[string]interface{}{
					"vpc": "vpc-test",
					"auth": map[string]interface{}{
						"authType": "secret",
						"secret": map[string]interface{}{
							"key":       "key.json",
							"name":      "key-secret",
							"namespace": "default",
						},
					},
				},
			},
			want: want{
				err:      errors.OK,
				provider: expectedProvider1,
			},
		}, {
			name: "FromMap with invalid data",
			fields: fields{
				Provider: provider,
			},
			args: args{
				m: map[string]interface{}{
					"vpc": "vpc-test",
					"auth": map[string]interface{}{
						"authType": "azaz",
					},
				},
			},
			want: want{
				err:      errors.InvalidArgument,
				provider: expectedProvider2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curProvider := tt.fields.Provider
			got := curProvider.FromMap(tt.args.m)
			if got.Code != tt.want.err.Code {
				t.Errorf("FromMap() = %v, want %v", got.Code, tt.want.err.Code)
			}
			if !curProvider.Equals(tt.want.provider) {
				t.Errorf("FromMap() = %v, want %v", curProvider, tt.want.provider)
			}
		})
	}
}

func TestProvider_Insert(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Provider
		KubernetesResources kubernetes.ResourceList
		Type                common.ProviderType
		Auth                auth.Auth
		VPCs                VPCCollection
	}
	type args struct {
		project Project
		update  []bool
	}
	type want struct {
		err      errors.Error
		provider Provider
	}
	testProject := NewProject("test", "owner_test")
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Insert non existing provider (creation)",
			fields: fields{
				Identifier: identifier.Provider{
					ID: "test",
				},
			},
			args: args{
				testProject,
				[]bool{},
			},
			want: want{
				err: errors.OK,
				provider: Provider{
					Identifier: identifier.Provider{
						ID: "test",
					},
				},
			},
		},
		{
			name: "Update existing provider (update)",
			fields: fields{
				Type: common.AWS,
				Identifier: identifier.Provider{
					ID: "test",
				},
			},
			args: args{
				testProject,
				[]bool{true},
			},
			want: want{
				err: errors.OK,
				provider: Provider{
					Type: common.AWS,
					Identifier: identifier.Provider{
						ID: "test",
					},
				},
			},
		},
		{
			name: "Update existing provider (no update)",
			fields: fields{
				Type: common.Azure,
				Identifier: identifier.Provider{
					ID: "test",
				},
			},
			args: args{
				testProject,
				[]bool{},
			},
			want: want{
				err: errors.Conflict,
				provider: Provider{
					Type: common.AWS,
					Identifier: identifier.Provider{
						ID: "test",
					},
				},
			},
		},
		{
			name: "Update non existing provider",
			fields: fields{
				Type: common.GCP,
				Identifier: identifier.Provider{
					ID: "test-2",
				},
			},
			args: args{
				testProject,
				[]bool{true},
			},
			want: want{
				err:      errors.NotFound,
				provider: Provider{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &Provider{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Type:                tt.fields.Type,
				Auth:                tt.fields.Auth,
				VPCs:                tt.fields.VPCs,
			}
			if got := provider.Insert(tt.args.project, tt.args.update...); got.Code != tt.want.err.Code {
				t.Errorf("Insert() = %v, want %v", got.Code, tt.want.err.Code)
			}
			if !reflect.DeepEqual(testProject.Resources[tt.fields.Identifier.ID], tt.want.provider) {
				t.Errorf("Insert() = %v, want %v", testProject.Resources[tt.fields.Identifier.ID], tt.want.provider)
			}
		})
	}
}

func TestProvider_Remove(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Provider
		KubernetesResources kubernetes.ResourceList
		Type                common.ProviderType
		Auth                auth.Auth
		VPCs                VPCCollection
	}
	type args struct {
		project Project
	}
	type want struct {
		err errors.Error
	}
	testProject := NewProject("test", "owner_test")
	providerTest1 := Provider{
		Identifier: identifier.Provider{
			ID: "test-1",
		},
	}
	testProject.Resources[providerTest1.Identifier.ID] = providerTest1
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Remove existing provider",
			fields: fields{
				Identifier: identifier.Provider{
					ID: "test-1",
				},
			},
			args: args{
				testProject,
			},
			want: want{
				err: errors.NoContent,
			},
		},
		{
			name: "Remove non-existing provider",
			fields: fields{
				Identifier: identifier.Provider{
					ID: "test-2",
				},
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
			provider := &Provider{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Type:                tt.fields.Type,
				Auth:                tt.fields.Auth,
				VPCs:                tt.fields.VPCs,
			}
			if got := provider.Remove(tt.args.project); got.Code != tt.want.err.Code {
				t.Errorf("Remove() = %v, want %v", got.Code, tt.want.err.Code)
			}
		})
	}
}
