package gcp

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upbound/provider-gcp/apis/compute/v1beta1"
	"testing"
)

type wantVM struct {
	Name   string
	Labels map[string]string
	Spec   v1beta1.InstanceSpec
}

func initVM(t *testing.T) (context.Context, *testResource.TestResource, usecase.VM) {
	ctx, testNet, suc := initSubnetwork(t)
	parentID := ctx.Value(testResource.NetworkIDKey).(identifier.Network)
	req := dto.CreateSubnetworkRequest{
		ParentID:    parentID,
		Name:        "test-subnet",
		Region:      "europe-west1",
		IPCIDRRange: "10.0.0.5/28",
		Managed:     false,
	}
	ctx.Set(context.RequestKey, req)
	subnet := &resource.Subnetwork{}
	err := suc.Create(ctx, subnet)
	require.True(t, err.IsOk())
	ctx.Set(testResource.SubnetworkIDKey, subnet.IdentifierID)
	uc := usecase.NewVMUseCase(testNet.ProjectRepo, repository.NewSSHKeyRepository(), gcp.NewRepository(repository.NewSSHKeyRepository()), nil, nil)
	return ctx, testNet, uc
}

func Test_vmUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, vuc := initVM(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clear(ctx)
	}()
	t.Run("Create a valid vm", func(t *testing.T) {
		req := dto.CreateVMRequest{
			ParentID:       ctx.Value(testResource.SubnetworkIDKey).(identifier.Subnetwork),
			Name:           "test-vm",
			AssignPublicIP: true,
			Zone:           "europe-west9-b",
			MachineType:    "e2-medium",
			Auths: []dto.VMAuth{
				{
					Username: "admin",
				}, {
					Username:     "test",
					RSAKeyLength: 1024,
				},
			},
			OS: resource.VMOS{
				Name: "ubuntu-1804-bionic-v20210223",
			},
			Managed: true,
		}
		ctx.Set(context.RequestKey, req)
		vm := &resource.VM{}
		err := vuc.Create(ctx, vm)
		assert.Equal(t, errors.Created, err)
	})

	t.Run("Create a vm with a non-existing secret should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a vm with an already existing name should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_vmUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
	}()
	t.Run("Delete a valid vm should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing vm should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a vm with children should fail", func(t *testing.T) {

	})
	t.Run("Delete cascade a vm should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
}
