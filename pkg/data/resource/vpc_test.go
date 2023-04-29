package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
	"testing"
)

func TestVPC_FromMap(t *testing.T) {
	type fields struct {
		VPC VPC
	}
	type args struct {
		data map[string]interface{}
	}
	type want struct {
		err errors.Error
		vpc VPC
	}
	vpc := NewVPC(identifier.VPC{ID: "test", ProviderID: "test"}, common.Azure)
	expectedVPC := vpc
	expectedVPC.Provider = "test"

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "FromMap with valid data",
			fields: fields{
				VPC: vpc,
			},
			args: args{
				data: map[string]interface{}{
					"provider": "test",
				},
			},
			want: want{
				err: errors.OK,
				vpc: expectedVPC,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curVpc := tt.fields.VPC
			if got := curVpc.FromMap(tt.args.data); !reflect.DeepEqual(got.Code, tt.want.err.Code) {
				t.Errorf("FromMap() = %v, want %v", got, tt.want)
			}
			if !curVpc.Equals(tt.want.vpc) {
				t.Errorf("FromMap() = %v, want %v", curVpc, tt.want.vpc)
			}
		})
	}
}

func TestVPC_Insert(t *testing.T) {
	type fields struct {
		vpc VPC
	}
	type args struct {
		project Project
		update  []bool
	}
	type want struct {
		err errors.Error
		vpc VPC
	}
	providerID := "test"
	testProject := NewProject("test", "owner_test")
	provider := NewProvider(identifier.Provider{ID: providerID}, common.AWS)
	testProject.Resources[providerID] = provider
	vpc1 := NewVPC(identifier.VPC{ID: "test", ProviderID: providerID}, common.AWS)
	vpc2 := NewVPC(identifier.VPC{ID: "test-2", ProviderID: providerID}, common.Azure)
	vpc3 := vpc1
	vpc3.Metadata.Tags = map[string]string{"test": "test"}
	vpc4 := vpc1
	vpc4.Metadata.Tags = map[string]string{"hello": "world"}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Insert non existing provider (creation)",
			fields: fields{
				vpc: vpc1,
			},
			args: args{
				testProject,
				[]bool{},
			},
			want: want{
				err: errors.OK,
				vpc: vpc1,
			},
		},
		{
			name: "Update existing provider (update)",
			fields: fields{
				vpc: vpc3,
			},
			args: args{
				testProject,
				[]bool{true},
			},
			want: want{
				err: errors.OK,
				vpc: vpc3,
			},
		},
		{
			name: "Update existing provider (no update)",
			fields: fields{
				vpc: vpc4,
			},
			args: args{
				testProject,
				[]bool{},
			},
			want: want{
				err: errors.Conflict,
				vpc: vpc3,
			},
		},
		{
			name: "Update non existing provider",
			fields: fields{
				vpc: vpc2,
			},
			args: args{
				testProject,
				[]bool{true},
			},
			want: want{
				err: errors.NotFound,
				vpc: VPC{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vpc := tt.fields.vpc
			if got := vpc.Insert(tt.args.project, tt.args.update...); got.Code != tt.want.err.Code {
				t.Errorf("Insert() = %v, want %v", got.Code, tt.want.err.Code)
			}
			id := tt.fields.vpc.Identifier
			vpcGot := testProject.Resources[id.ProviderID].VPCs[id.ID]
			if !vpcGot.Equals(tt.want.vpc) {
				t.Errorf("Insert() = %v, want %v", testProject.Resources[providerID].VPCs[tt.fields.vpc.Identifier.ID], tt.want.vpc)
			}
		})
	}
}
