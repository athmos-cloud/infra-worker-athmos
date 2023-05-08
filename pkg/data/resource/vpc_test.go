package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
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
	vpc := NewVPC(NewResourcePayload{
		Name: "test",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: "test",
		}),
		Provider: types.GCP,
	})
	expectedVPC := vpc

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
				vpc: expectedVPC,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curVpc := tt.fields.VPC
			defer func() {
				if r := recover(); r != nil {
					err := r.(errors.Error)
					if err.Code != tt.want.err.Code {
						t.Errorf("FromMap()  %v, want %v", err.Code, tt.want.err.Code)
					}
				}
			}()
			curVpc.FromMap(tt.args.data)
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
		update  bool
	}
	type want struct {
		err errors.Error
		vpc VPC
	}
	providerID := "test"
	testProject := NewProject("test", "owner_test")
	provider := NewProvider(NewResourcePayload{
		Name: "test",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: "test",
		}),
		Provider: types.GCP,
	})
	testProject.Resources[providerID] = provider
	vpc1 := NewVPC(NewResourcePayload{
		Name: "test",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: "test",
		}),
		Provider: types.GCP,
	})
	vpc1.Identifier.VPCID = "test"
	vpc2 := NewVPC(NewResourcePayload{
		Name: "test2",
		ParentIdentifier: identifier.Build(identifier.IdPayload{
			ProviderID: "test",
		}),
		Provider: types.GCP,
	})
	vpc2.Identifier.VPCID = "test2"
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
				*testProject,
				false,
			},
			want: want{
				vpc: vpc1,
			},
		},
		{
			name: "Update existing provider (update)",
			fields: fields{
				vpc: vpc3,
			},
			args: args{
				*testProject,
				true,
			},
			want: want{
				vpc: vpc3,
			},
		},
		{
			name: "Update existing provider (no update)",
			fields: fields{
				vpc: vpc4,
			},
			args: args{
				*testProject,
				false,
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
				*testProject,
				true,
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
			defer func() {
				if r := recover(); r != nil {
					err := r.(errors.Error)
					if err.Code != tt.want.err.Code {
						t.Errorf("Insert()  %v, want %v", err.Code, tt.want.err.Code)
					}
				}
			}()
			if tt.args.update {
				tt.args.project.Update(&vpc)
			} else {
				tt.args.project.Insert(&vpc)
			}
			id := tt.fields.vpc.Identifier
			vpcGot := tt.args.project.Resources[id.ProviderID].VPCs[id.VPCID]
			if !vpcGot.Equals(tt.want.vpc) {
				t.Errorf("Insert() = %v, want %v", testProject.Resources[providerID].VPCs[tt.fields.vpc.Identifier.VPCID], tt.want.vpc)
			}
		})
	}
}
