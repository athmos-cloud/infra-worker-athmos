package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upbound/provider-aws/apis/ec2/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"strings"
	"testing"
)

const (
	testSubnetName = "test-subnet"
)

type wantVM struct {
	Name   string
	Labels map[string]string
	Spec   v1beta1.InstanceSpec
}

func vmSuiteTeardown(ctx context.Context, t *testing.T, container *gnomock.Container) {
	require.NoError(t, gnomock.Stop(container))
	ClearFixtures(ctx)
}

func vmTeardown(ctx context.Context) {
	ClearVMFixtures(ctx)
}

func Test_vmUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	sshRepo := repository.NewSSHKeyRepository()

	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	vuc := usecase.NewVMUseCase(resourceTest.ProjectRepo, sshRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)
	SubnetworkFixture(ctx, t, suc)

	defer vmSuiteTeardown(ctx, t, mongoC)

	t.Run("Create a valid vm", func(t *testing.T) {
		defer vmTeardown(ctx)

		vm := VMFixture(ctx, t, vuc)

		kubeResource := &v1beta1.Instance{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: vm.IdentifierID.VM}, kubeResource)
		require.NoError(t, errk)
		wantLabels := map[string]string{
			"app.kubernetes.io/managed-by": "athmos",
			"athmos.cloud/project-id":      ctx.Value(context.ProjectIDKey).(string),
			"identifier.provider":          vm.IdentifierID.Provider,
			"identifier.vpc":               vm.IdentifierID.VPC,
			"identifier.network":           vm.IdentifierID.Network,
			"identifier.subnetwork":        vm.IdentifierID.Subnetwork,
			"identifier.vm":                vm.IdentifierID.VM,
			"vm-ssh-keys_admin":            vm.Auths[0].SecretName,
			"vm-ssh-keys_test":             vm.Auths[1].SecretName,
			"name.provider":                "fixture-provider",
			"name.vpc":                     "test",
			"name.network":                 "fixture-network",
			"name.subnetwork":              "fixture-subnet",
			"name.vm":                      vm.IdentifierName.VM,
			"vm-has-public-ip":             "true",
			"vm-ssh-keys-secret-namespace": ctx.Value(test.TestNamespaceContextKey).(string),
		}

		keyName := fmt.Sprintf("%s-keypair", vm.IdentifierID.VM)
		sizeGib := float64(10)

		kind := "instance.ec2.aws.upbound.io"
		tags := map[string]*string{
			"crossplane-kind":           &kind,
			"crossplane-name":           &vm.IdentifierID.VM,
			"crossplane-providerconfig": &vm.IdentifierID.VM,
		}

		diskType := "gp2"
		sshKeys := ""
		for _, auth := range vm.Auths {
			sshKeys += fmt.Sprintf("%s:%s\n", auth.Username, auth.PublicKey)
		}
		sshKeys = strings.TrimSuffix(sshKeys, "\n")
		subnetID := identifier.Subnetwork{
			Provider:   vm.IdentifierID.Provider,
			VPC:        vm.IdentifierID.VPC,
			Network:    vm.IdentifierID.Network,
			Subnetwork: vm.IdentifierID.Subnetwork,
		}

		wantSpec := v1beta1.InstanceSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: "Delete",
				ProviderConfigReference: &v1.Reference{
					Name: vm.IdentifierID.VM,
				},
				ManagementPolicy: "FullControl",
			},
			ForProvider: v1beta1.InstanceParameters{
				AMI:                      &vm.OS.ID,
				AssociatePublicIPAddress: &vm.AssignPublicIP,
				InstanceType:             &vm.MachineType,
				KeyName:                  &keyName,
				Region:                   &vm.Zone,
				RootBlockDevice: []v1beta1.RootBlockDeviceParameters{
					{
						DeleteOnTermination: &vm.Disks[0].AutoDelete,
						VolumeSize:          &sizeGib,
						VolumeType:          &diskType,
					},
				},
				SubnetIDSelector: &v1.Selector{
					MatchLabels: subnetID.ToIDLabels(),
				},
				Tags: tags,
			},
		}
		wantNet := wantVM{
			Name:   vm.IdentifierID.VM,
			Labels: wantLabels,
			Spec:   wantSpec,
		}
		gotNet := wantVM{
			Name:   kubeResource.Name,
			Labels: kubeResource.Labels,
			Spec:   kubeResource.Spec,
		}
		assert.Equal(t, wantNet, gotNet)
	})

	t.Run("Create a vm with an HDD root block should fail", func(t *testing.T) {
		defer vmTeardown(ctx)

		machineType := "t2.micro"
		zone := "eu-west-1"
		osName := "ami-0a5d9cd4e632d99c1"
		autoDelete := true
		req := dto.CreateVMRequest{
			ParentID:       ctx.Value(testResource.SubnetworkIDKey).(identifier.Subnetwork),
			Name:           "fixture-vm",
			AssignPublicIP: true,
			Zone:           zone,
			MachineType:    machineType,
			Auths: []dto.VMAuth{
				{
					Username: "admin",
				}, {
					Username:     "test",
					RSAKeyLength: 1024,
				},
			},
			OS: instance.VMOS{
				ID:   osName,
				Name: osName,
			},
			Disks: []instance.VMDisk{
				{
					AutoDelete: autoDelete,
					Mode:       instance.DiskModeReadWrite,
					Type:       instance.DiskTypeHDD,
					SizeGib:    10,
				},
			},
			Managed: true,
		}
		ctx.Set(context.RequestKey, req)

		vm := &instance.VM{}
		err := vuc.Create(ctx, vm)

		assert.Equal(t, errors.BadRequest.Code, err.Code)
	})

	t.Run("Create a vm with an already existing name should fail", func(t *testing.T) {
		defer vmTeardown(ctx)

		vm := VMFixture(ctx, t, vuc)
		ctx.Set(context.RequestKey, dto.CreateVMRequest{
			Name: vm.IdentifierName.VM,
			ParentID: identifier.Subnetwork{
				Provider:   vm.IdentifierID.Provider,
				VPC:        vm.IdentifierID.VPC,
				Network:    vm.IdentifierID.Network,
				Subnetwork: vm.IdentifierID.Subnetwork,
			},
		})
		toCreate := &instance.VM{}
		err := vuc.Create(ctx, toCreate)
		require.Equal(t, errors.Conflict.Code, err.Code)
	})
}

func Test_vmUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	sshRepo := repository.NewSSHKeyRepository()

	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	vuc := usecase.NewVMUseCase(resourceTest.ProjectRepo, sshRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)
	SubnetworkFixture(ctx, t, suc)

	defer vmSuiteTeardown(ctx, t, mongoC)

	t.Run("Delete a valid vm should succeed", func(t *testing.T) {
		defer vmTeardown(ctx)

		vm := VMFixture(ctx, t, vuc)
		ctx.Set(context.RequestKey, dto.DeleteVMRequest{
			IdentifierID: vm.IdentifierID,
		})
		foundVM := &instance.VM{}
		err := vuc.Delete(ctx, foundVM)
		assert.Equal(t, errors.NoContent.Code, err.Code)
	})
	t.Run("Delete a non-existing vm should return not found", func(t *testing.T) {
		defer vmTeardown(ctx)
		ctx.Set(context.RequestKey, dto.DeleteVMRequest{
			IdentifierID: identifier.VM{
				Provider:   "provider-test",
				Network:    "network-test",
				Subnetwork: "subnet-test",
				VM:         "this-vm-does-not-exist",
			},
		})
		foundVM := &instance.VM{}
		err := vuc.Delete(ctx, foundVM)
		assert.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func Test_vmUseCase_Get(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	sshRepo := repository.NewSSHKeyRepository()

	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	vuc := usecase.NewVMUseCase(resourceTest.ProjectRepo, sshRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)
	SubnetworkFixture(ctx, t, suc)

	defer vmSuiteTeardown(ctx, t, mongoC)

	t.Run("Get a valid vm should succeed", func(t *testing.T) {
		defer vmTeardown(ctx)

		vm := VMFixture(ctx, t, vuc)
		ctx.Set(context.RequestKey, dto.GetVMRequest{
			IdentifierID: vm.IdentifierID,
		})
		foundVM := &instance.VM{}
		err := vuc.Get(ctx, foundVM)
		vm.Metadata.Tags = map[string]string{}
		for _, auth := range vm.Auths {
			auth.KeyLength = 0
		}
		assert.Equal(t, errors.OK.Code, err.Code)
		assert.Equal(t, vm, foundVM)
	})
	t.Run("Get a non-existing vm should return not found", func(t *testing.T) {
		defer vmTeardown(ctx)

		ctx.Set(context.RequestKey, dto.GetVMRequest{
			IdentifierID: identifier.VM{
				Provider:   "provider-test",
				Network:    "network-test",
				Subnetwork: "subnet-test",
				VM:         "this-vm-does-not-exist",
			},
		})
		delVM := &instance.VM{}
		err := vuc.Get(ctx, delVM)
		assert.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func Test_vmUseCase_Update(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	sshRepo := repository.NewSSHKeyRepository()

	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	vuc := usecase.NewVMUseCase(resourceTest.ProjectRepo, sshRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)
	SubnetworkFixture(ctx, t, suc)

	defer vmSuiteTeardown(ctx, t, mongoC)

	t.Run("Update an existing vm should succeed", func(t *testing.T) {
		defer vmTeardown(ctx)

		vm := VMFixture(ctx, t, vuc)
		ctx.Set(context.RequestKey, dto.UpdateVMRequest{
			IdentifierID: vm.IdentifierID,
			Auths: &[]dto.VMAuth{
				{
					Username: "admin",
				},
				{
					Username: "some-new-user",
				},
			},
		})
		updatedVM := &instance.VM{}
		err := vuc.Update(ctx, updatedVM)
		assert.Equal(t, errors.NoContent.Code, err.Code)
		expectedUserList := []string{"admin", "some-new-user"}
		for _, auth := range updatedVM.Auths {
			assert.Contains(t, expectedUserList, auth.Username)
		}
	})
	t.Run("Update a non-existing VM should return not found error", func(t *testing.T) {
		defer vmTeardown(ctx)

		ctx.Set(context.RequestKey, dto.UpdateVMRequest{
			IdentifierID: identifier.VM{
				Provider:   "provider-test",
				Network:    "network-test",
				Subnetwork: "subnet-test",
				VM:         "this-vm-does-not-exist",
			},
		})
		delVM := &instance.VM{}
		err := vuc.Update(ctx, delVM)
		assert.Equal(t, errors.NotFound.Code, err.Code)
	})
}
