package gcp

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upbound/provider-gcp/apis/compute/v1beta1"
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

func initVM(t *testing.T) (context.Context, *testResource.TestResource, usecase.VM) {
	ctx, testNet, suc := initSubnetwork(t)
	parentID := ctx.Value(testResource.NetworkIDKey).(identifier.Network)
	req := dto.CreateSubnetworkRequest{
		ParentID:    parentID,
		Name:        testSubnetName,
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

func createVM(t *testing.T, ctx context.Context, vuc usecase.VM) *resource.VM {
	machineType := "e2-medium"
	zone := "europe-west9-b"
	osName := "ubuntu-1804-bionic-v20210223"
	autoDelete := true

	req := dto.CreateVMRequest{
		ParentID:       ctx.Value(testResource.SubnetworkIDKey).(identifier.Subnetwork),
		Name:           "test-vm",
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
		OS: resource.VMOS{
			ID: osName,
		},
		Disks: []resource.VMDisk{
			{
				AutoDelete: autoDelete,
				Mode:       resource.DiskModeReadWrite,
				Type:       resource.DiskTypeHDD,
				SizeGib:    10,
			},
		},
		Managed: true,
	}
	ctx.Set(context.RequestKey, req)
	vm := &resource.VM{}
	err := vuc.Create(ctx, vm)
	require.Equal(t, errors.Created, err)

	return vm
}

func Test_vmUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, vuc := initVM(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clear(ctx)
	}()
	t.Run("Create a valid vm", func(t *testing.T) {
		machineType := "e2-medium"
		zone := "europe-west9-b"
		osName := "ubuntu-1804-bionic-v20210223"
		autoDelete := true
		vm := createVM(t, ctx, vuc)
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
			"name.provider":                "test",
			"name.vpc":                     "test",
			"name.network":                 "test-net",
			"name.subnetwork":              testSubnetName,
			"name.vm":                      "test-vm",
			"vm-has-public-ip":             "true",
			"vm-ssh-keys-secret-namespace": ctx.Value(test.TestNamespaceContextKey).(string),
			"vm-ssh-keys-names":            fmt.Sprintf("%s.%s", vm.Auths[0].SecretName, vm.Auths[1].SecretName),
		}
		readWrite := "READ_WRITE"
		sizeGib := float64(10)
		diskType := "pd-standard"
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
		netID := identifier.Network{
			Provider: vm.IdentifierID.Provider,
			VPC:      vm.IdentifierID.VPC,
			Network:  vm.IdentifierID.Network,
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
				Project:     &vm.IdentifierID.VPC,
				MachineType: &machineType,
				Zone:        &zone,
				BootDisk: []v1beta1.BootDiskParameters{
					{
						AutoDelete: &autoDelete,
						Mode:       &readWrite,
						InitializeParams: []v1beta1.InitializeParamsParameters{
							{
								Image: &osName,
								Type:  &diskType,
								Size:  &sizeGib,
							},
						},
					},
				},
				Metadata: map[string]*string{
					"ssh-keys": &sshKeys,
				},
				NetworkInterface: []v1beta1.NetworkInterfaceParameters{
					{
						NetworkSelector: &v1.Selector{
							MatchLabels: netID.ToIDLabels(),
						},
						SubnetworkSelector: &v1.Selector{
							MatchLabels: subnetID.ToIDLabels(),
						},
						AccessConfig: []v1beta1.AccessConfigParameters{
							{
								NATIP:               nil,
								NetworkTier:         nil,
								PublicPtrDomainName: nil,
							},
						},
					},
				},
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

	t.Run("Create a vm with an already existing name should fail", func(t *testing.T) {
		vm := createVM(t, ctx, vuc)
		ctx.Set(context.RequestKey, dto.CreateVMRequest{
			Name: vm.IdentifierName.VM,
			ParentID: identifier.Subnetwork{
				Provider:   vm.IdentifierID.Provider,
				VPC:        vm.IdentifierID.VPC,
				Network:    vm.IdentifierID.Network,
				Subnetwork: vm.IdentifierID.Subnetwork,
			},
		})
		toCreate := &resource.VM{}
		err := vuc.Create(ctx, toCreate)
		require.Equal(t, errors.Conflict.Code, err.Code)
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
