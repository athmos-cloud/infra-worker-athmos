package helm

import (
	"context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	helmclient "github.com/mittwald/go-helm-client"
	"reflect"
	"testing"
)

func TestReleaseRepository_Create(t *testing.T) {
	type fields struct {
		HelmClient helmclient.Client
		Namespace  string
	}
	type args struct {
		ctx     context.Context
		request option.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReleaseDAO{
				HelmClient: tt.fields.HelmClient,
				Namespace:  tt.fields.Namespace,
			}
			got, got1 := r.Create(tt.args.ctx, tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateProject() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CreateProject() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestReleaseRepository_Delete(t *testing.T) {
	type fields struct {
		HelmClient helmclient.Client
		Namespace  string
	}
	type args struct {
		ctx  context.Context
		optn option.Option
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
			r := &ReleaseDAO{
				HelmClient: tt.fields.HelmClient,
				Namespace:  tt.fields.Namespace,
			}
			if got := r.Delete(tt.args.ctx, tt.args.optn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReleaseRepository_Get(t *testing.T) {
	type fields struct {
		HelmClient helmclient.Client
		Namespace  string
	}
	type args struct {
		in0  context.Context
		optn option.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReleaseDAO{
				HelmClient: tt.fields.HelmClient,
				Namespace:  tt.fields.Namespace,
			}
			got, got1 := r.Get(tt.args.in0, tt.args.optn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestReleaseRepository_GetAll(t *testing.T) {
	type fields struct {
		HelmClient helmclient.Client
		Namespace  string
	}
	type args struct {
		ctx_ context.Context
		in1  option.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []interface{}
		want1  errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReleaseDAO{
				HelmClient: tt.fields.HelmClient,
				Namespace:  tt.fields.Namespace,
			}
			got, got1 := r.GetAll(tt.args.ctx_, tt.args.in1)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetAll() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestReleaseRepository_Update(t *testing.T) {
	type fields struct {
		HelmClient helmclient.Client
		Namespace  string
	}
	type args struct {
		ctx  context.Context
		optn option.Option
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
			r := &ReleaseDAO{
				HelmClient: tt.fields.HelmClient,
				Namespace:  tt.fields.Namespace,
			}
			if got := r.Update(tt.args.ctx, tt.args.optn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateProjectName() = %v, want %v", got, tt.want)
			}
		})
	}
}
