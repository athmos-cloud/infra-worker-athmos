package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	networkModel "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
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
	"testing"
)

type wantNetwork struct {
	Name   string
	Labels map[string]string
	Spec   v1beta1.VPCSpec
}

func netSuiteTeardown(ctx context.Context, t *testing.T, container *gnomock.Container) {
	require.NoError(t, gnomock.Stop(container))
	ClearFixtures(ctx)
}

func netTeardown(ctx context.Context) {
	ClearNetworksFixtures(ctx)
}

func Test_networkUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(
		resourceTest.ProjectRepo,
		nil,
		awsRepo,
		nil)

	ProviderFixture(ctx, t, puc)

	defer netSuiteTeardown(ctx, t, mongoC)

	t.Run("Create a valid network", func(t *testing.T) {
		defer netTeardown(ctx)

		net := NetworkFixture(ctx, t, nuc)

		kubeResource := &v1beta1.VPC{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: net.IdentifierID.Network}, kubeResource)
		assert.NoError(t, errk)
		wantLabels := map[string]string{
			"app.kubernetes.io/managed-by": "athmos",
			"athmos.cloud/project-id":      ctx.Value(context.ProjectIDKey).(string),
			"identifier.provider":          net.IdentifierID.Provider,
			"identifier.vpc":               net.IdentifierID.VPC,
			"identifier.network":           net.IdentifierID.Network,
			"name.provider":                "fixture-provider",
			"name.vpc":                     "test",
			"name.network":                 net.IdentifierName.Network,
		}
		//autoCreateSubnet := false
		region := "eu-west-1"

		wantSpec := v1beta1.VPCSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: "Delete",
				ProviderConfigReference: &v1.Reference{
					Name: net.IdentifierID.Provider,
				},
				ManagementPolicy: "FullControl",
			},
			ForProvider: v1beta1.VPCParameters_2{
				Region: &region,
			},
		}
		wantNet := wantNetwork{
			Name:   net.IdentifierID.Network,
			Labels: wantLabels,
			Spec:   wantSpec,
		}
		gotNet := wantNetwork{
			Name:   kubeResource.Name,
			Labels: kubeResource.Labels,
			Spec:   kubeResource.Spec,
		}
		assertNetworkEqual(t, wantNet, gotNet)
	})

	t.Run("Create a network without specifying a region should raise a bad request error", func(t *testing.T) {
		defer netTeardown(ctx)

		parentID := ctx.Value(testResource.ProviderIDKey).(identifier.Provider)
		netName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))
		req := dto.CreateNetworkRequest{
			ParentIDProvider: &parentID,
			Name:             netName,
			Managed:          false,
		}
		ctx.Set(context.RequestKey, req)

		network := &networkModel.Network{}
		err := nuc.Create(ctx, network)

		ctx.Set(testResource.NetworkIDKey, network.IdentifierID)

		require.Equal(t, errors.BadRequest.Code, err.Code)
	})

	t.Run("Create a network with an already existing name should raise conflict error", func(t *testing.T) {
		defer netTeardown(ctx)

		net := NetworkFixture(ctx, t, nuc)
		err := nuc.Create(ctx, net)
		require.Equal(t, errors.Conflict.Code, err.Code)
	})
}

func Test_networkUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(
		resourceTest.ProjectRepo,
		nil,
		awsRepo,
		nil)

	ProviderFixture(ctx, t, puc)

	defer netSuiteTeardown(ctx, t, mongoC)

	t.Run("Delete a valid network should succeed", func(t *testing.T) {
		defer netTeardown(ctx)

		net := NetworkFixture(ctx, t, nuc)
		delReq := dto.DeleteNetworkRequest{
			IdentifierID: net.IdentifierID,
		}
		ctx.Set(context.RequestKey, delReq)
		err := nuc.Delete(ctx, net)
		require.Equal(t, errors.NoContent.Code, err.Code)
	})

	t.Run("Delete a non-existing network should fail", func(t *testing.T) {
		defer netTeardown(ctx)

		delReq := dto.DeleteNetworkRequest{
			IdentifierID: identifier.Network{
				Provider: "test",
				VPC:      "test",
				Network:  "this-network-does-not-exist",
			},
		}
		ctx.Set(context.RequestKey, delReq)
		net := &network.Network{}
		err := nuc.Delete(ctx, net)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})

	t.Run("Delete a network with children should cascade", func(t *testing.T) {
		defer netTeardown(ctx)

		sshRepo := repository.NewSSHKeyRepository()
		suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
		vuc := usecase.NewVMUseCase(resourceTest.ProjectRepo, sshRepo, nil, awsRepo, nil)

		network := NetworkFixture(ctx, t, nuc)
		SubnetworkFixture(ctx, t, suc)
		VMFixture(ctx, t, vuc)

		delReq := dto.DeleteNetworkRequest{
			IdentifierID: network.IdentifierID,
		}
		ctx.Set(context.RequestKey, delReq)
		err := nuc.Delete(ctx, network)

		require.Equal(t, errors.NoContent.Code, err.Code)
	})
}

func Test_networkUseCase_Get(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(
		resourceTest.ProjectRepo,
		nil,
		awsRepo,
		nil)

	ProviderFixture(ctx, t, puc)

	defer netSuiteTeardown(ctx, t, mongoC)

	t.Run("Get a valid network should succeed", func(t *testing.T) {
		defer netTeardown(ctx)

		net := NetworkFixture(ctx, t, nuc)
		getReq := dto.GetNetworkRequest{
			IdentifierID: net.IdentifierID,
		}
		ctx.Set(context.RequestKey, getReq)
		gotNet := &network.Network{}
		err := nuc.Get(ctx, gotNet)
		require.Equal(t, errors.OK.Code, err.Code)
		assert.Equal(t, net, gotNet)
	})

	t.Run("Delete a non-existing network should fail", func(t *testing.T) {
		defer netTeardown(ctx)

		getReq := dto.GetNetworkRequest{
			IdentifierID: identifier.Network{
				Provider: "test",
				VPC:      "test",
				Network:  "this-network-does-not-exist",
			},
		}
		ctx.Set(context.RequestKey, getReq)
		gotNet := &network.Network{}
		err := nuc.Get(ctx, gotNet)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})

}

func Test_networkUseCase_Update(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(
		resourceTest.ProjectRepo,
		nil,
		awsRepo,
		nil)

	ProviderFixture(ctx, t, puc)

	defer netSuiteTeardown(ctx, t, mongoC)

	t.Run("Update a valid network should succeed", func(t *testing.T) {
		defer netTeardown(ctx)

		net := NetworkFixture(ctx, t, nuc)
		managed := false
		toUp := dto.UpdateNetworkRequest{
			IdentifierID: net.IdentifierID,
			Managed:      &managed,
		}
		ctx.Set(context.RequestKey, toUp)
		err := nuc.Update(ctx, net)
		//fmt.Println(err.Message)
		require.Equal(t, errors.NoContent.Code, err.Code)

		kubeResource := &v1beta1.VPC{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: net.IdentifierID.Network}, kubeResource)

		assert.NoError(t, errk)
		assert.Equal(t, kubeResource.Spec.DeletionPolicy, v1.DeletionOrphan)
	})

	t.Run("Update a non-existing network should fail", func(t *testing.T) {
		defer netTeardown(ctx)

		toUp := dto.UpdateNetworkRequest{
			IdentifierID: identifier.Network{
				Provider: "test",
				VPC:      "test",
				Network:  "this-network-does-not-exist",
			},
		}
		ctx.Set(context.RequestKey, toUp)
		net := &network.Network{}
		err := nuc.Update(ctx, net)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func assertNetworkEqual(t *testing.T, want wantNetwork, got wantNetwork) {
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Labels, got.Labels)
	assert.Equal(t, want.Spec, got.Spec)
}
