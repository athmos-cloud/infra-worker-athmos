package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	domain "github.com/athmos-cloud/infra-worker-athmos/pkg/data/auth"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
	"testing"
)

func TestProvider_FromMap(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Provider
		KubernetesResources kubernetes.ResourceList
		Type                common.ProviderType
		Auth                domain.Auth
		VPCs                VPCCollection
	}
	type args struct {
		m map[string]interface{}
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
			provider := &Provider{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Type:                tt.fields.Type,
				Auth:                tt.fields.Auth,
				VPCs:                tt.fields.VPCs,
			}
			if got := provider.FromMap(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromMap() = %v, want %v", got, tt.want)
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
		Auth                domain.Auth
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

func TestProvider_ToDomain(t *testing.T) {
	type fields struct {
		Metadata            metadata.Metadata
		Identifier          identifier.Provider
		KubernetesResources kubernetes.ResourceList
		Type                common.ProviderType
		Auth                domain.Auth
		VPCs                VPCCollection
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
			provider := &Provider{
				Metadata:            tt.fields.Metadata,
				Identifier:          tt.fields.Identifier,
				KubernetesResources: tt.fields.KubernetesResources,
				Type:                tt.fields.Type,
				Auth:                tt.fields.Auth,
				VPCs:                tt.fields.VPCs,
			}
			got, got1 := provider.ToDomain()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDomain() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ToDomain() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
