package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upbound/provider-aws/apis/ec2/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

type wantSubnetwork struct {
	Name   string
	Labels map[string]string
	Spec   v1beta1.SubnetSpec
}

func subSuiteTeardown(ctx context.Context, t *testing.T, container *gnomock.Container) {
	require.NoError(t, gnomock.Stop(container))
	ClearFixtures(ctx)
}

func subTeardown(ctx context.Context) {
	ClearSubnetworkFixtures(ctx)
}

func Test_subnetworkUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)

	defer subSuiteTeardown(ctx, t, mongoC)

	t.Run("Create a valid subnetwork", func(t *testing.T) {
		defer subTeardown(ctx)

		subnet := SubnetworkFixture(ctx, t, suc)
		region := "eu-west-1"
		ipCIDR := "10.0.0.1/26"

		kubeResource := &v1beta1.Subnet{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: subnet.IdentifierID.Subnetwork}, kubeResource)
		assert.NoError(t, errk)

		wantLabels := map[string]string{
			"app.kubernetes.io/managed-by": "athmos",
			"athmos.cloud/project-id":      ctx.Value(context.ProjectIDKey).(string),
			"identifier.provider":          subnet.IdentifierID.Provider,
			"identifier.vpc":               subnet.IdentifierID.VPC,
			"identifier.network":           subnet.IdentifierID.Network,
			"identifier.subnetwork":        subnet.IdentifierID.Subnetwork,
			"name.provider":                "fixture-provider",
			"name.vpc":                     "test",
			"name.network":                 "fixture-network",
			"name.subnetwork":              subnet.IdentifierName.Subnetwork,
		}
		wantSpec := v1beta1.SubnetSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: "Delete",
				ProviderConfigReference: &v1.Reference{
					Name: subnet.IdentifierID.Provider,
				},
				ManagementPolicy: "FullControl",
			},
			ForProvider: v1beta1.SubnetParameters_2{
				VPCID:     &subnet.IdentifierID.Network,
				Region:    &region,
				CidrBlock: &ipCIDR,
			},
		}
		wantNet := wantSubnetwork{
			Name:   subnet.IdentifierID.Subnetwork,
			Labels: wantLabels,
			Spec:   wantSpec,
		}
		gotNet := wantSubnetwork{
			Name:   kubeResource.Name,
			Labels: kubeResource.Labels,
			Spec:   kubeResource.Spec,
		}
		assertSubnetworkEqual(t, wantNet, gotNet)
	})

	t.Run("Create a subnetwork with an already existing name should return conflict error", func(t *testing.T) {
		_ = SubnetworkFixture(ctx, t, suc)
		newSubnet := &network.Subnetwork{}
		err := suc.Create(ctx, newSubnet)
		require.Equal(t, errors.Conflict.Code, err.Code)
	})
}

func Test_subnetworkUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)

	defer subSuiteTeardown(ctx, t, mongoC)

	t.Run("Delete a valid subnetwork should succeed", func(t *testing.T) {
		subnet := SubnetworkFixture(ctx, t, suc)
		delReq := dto.DeleteSubnetworkRequest{
			IdentifierID: subnet.IdentifierID,
		}
		ctx.Set(context.RequestKey, delReq)
		delSubnet := &network.Subnetwork{}
		err := suc.Delete(ctx, delSubnet)
		assert.Equal(t, errors.NoContent.Code, err.Code)
	})
	t.Run("Delete a non-existing subnetwork should fail", func(t *testing.T) {
		delReq := dto.DeleteSubnetworkRequest{
			IdentifierID: identifier.Subnetwork{
				Provider:   "test",
				VPC:        "test",
				Network:    "test",
				Subnetwork: "this-network-does-not-exist",
			},
		}
		ctx.Set(context.RequestKey, delReq)
		subnet := &network.Subnetwork{}
		err := suc.Delete(ctx, subnet)
		assert.Equal(t, errors.NotFound.Code, err.Code)
	})
	t.Run("Delete a subnetwork with children should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete cascade a subnetwork should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_subnetworkUseCase_Get(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)

	defer subSuiteTeardown(ctx, t, mongoC)

	t.Run("Get a valid subnetwork should succeed", func(t *testing.T) {
		parentID := ctx.Value(testResource.NetworkIDKey).(identifier.Network)
		subnetName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))
		region := "europe-west1"
		ipCIDR := "10.0.0.1/26"
		req := dto.CreateSubnetworkRequest{
			ParentID:    parentID,
			Name:        subnetName,
			Region:      region,
			IPCIDRRange: ipCIDR,
			Managed:     false,
		}
		ctx.Set(context.RequestKey, req)
		subnet := &network.Subnetwork{}
		err := suc.Create(ctx, subnet)
		assert.Equal(t, errors.Created.Code, err.Code)
		getReq := dto.GetSubnetworkRequest{
			IdentifierID: subnet.IdentifierID,
		}
		ctx.Set(context.RequestKey, getReq)
		getSubnet := &network.Subnetwork{}
		err = suc.Get(ctx, getSubnet)
		assert.Equal(t, errors.OK.Code, err.Code)
		assert.Equal(t, subnet.IdentifierID, getReq.IdentifierID)
		assert.Equal(t, subnet.IPCIDRRange, ipCIDR)
		assert.Equal(t, subnet.Region, region)
	})
	t.Run("Get a non-existing subnetwork should fail", func(t *testing.T) {
		getReq := dto.GetSubnetworkRequest{
			IdentifierID: identifier.Subnetwork{
				Provider:   "test",
				VPC:        "test",
				Network:    "test",
				Subnetwork: "this-network-does-not-exist",
			},
		}
		ctx.Set(context.RequestKey, getReq)
		subnet := &network.Subnetwork{}
		err := suc.Get(ctx, subnet)
		assert.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func Test_subnetworkUseCase_Update(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)

	defer subSuiteTeardown(ctx, t, mongoC)

	t.Run("Update a valid subnetwork should succeed", func(t *testing.T) {
		parentID := ctx.Value(testResource.NetworkIDKey).(identifier.Network)
		subnetName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))
		region := "eu-west-1"
		ipCIDR := "10.0.0.1/26"
		req := dto.CreateSubnetworkRequest{
			ParentID:    parentID,
			Name:        subnetName,
			Region:      region,
			IPCIDRRange: ipCIDR,
			Managed:     false,
		}
		ctx.Set(context.RequestKey, req)
		subnet := &network.Subnetwork{}
		err := suc.Create(ctx, subnet)
		assert.Equal(t, errors.Created.Code, err.Code)
		newRegion := "eu-west-2"
		newIPCIDR := "10.1.0.1/26"
		updReq := dto.UpdateSubnetworkRequest{
			IdentifierID: subnet.IdentifierID,
			Region:       &newRegion,
			IPCIDRRange:  &newIPCIDR,
		}
		ctx.Set(context.RequestKey, updReq)
		err = suc.Update(ctx, subnet)
		assert.Equal(t, errors.NoContent.Code, err.Code)

		kubeSubnet := &v1beta1.Subnet{}
		errKube := kubernetes.Client().Client.Get(ctx, client.ObjectKey{Name: subnet.IdentifierID.Subnetwork}, kubeSubnet)
		assert.NoError(t, errKube)
		assert.Equal(t, newRegion, *kubeSubnet.Spec.ForProvider.Region)
		assert.Equal(t, newIPCIDR, *kubeSubnet.Spec.ForProvider.CidrBlock)
	})
	t.Run("Update a non-existing subnetwork should fail", func(t *testing.T) {
		newRegion := "europe-west2"
		newIPCIDR := "10.1.0.1/26"
		updReq := dto.UpdateSubnetworkRequest{
			IdentifierID: identifier.Subnetwork{
				Provider:   "test",
				VPC:        "test",
				Network:    "test",
				Subnetwork: "this-network-does-not-exist",
			},
			Region:      &newRegion,
			IPCIDRRange: &newIPCIDR,
		}
		ctx.Set(context.RequestKey, updReq)
		subnet := &network.Subnetwork{}
		err := suc.Update(ctx, subnet)
		assert.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func assertSubnetworkEqual(t *testing.T, want wantSubnetwork, got wantSubnetwork) {
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Labels, got.Labels)
	assert.Equal(t, want.Spec, got.Spec)
}
