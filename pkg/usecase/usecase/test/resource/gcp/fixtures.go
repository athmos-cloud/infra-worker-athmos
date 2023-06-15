package gcp

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func ProviderFixture(ctx context.Context, t *testing.T, puc usecase.Provider) *resource.Provider {
	req := dto.CreateProviderRequest{
		Name:           "test",
		VPC:            testResource.SecretTestName,
		SecretAuthName: testResource.SecretTestName,
	}
	ctx.Set(context.RequestKey, req)

	provider := &resource.Provider{}
	err := puc.Create(ctx, provider)
	assert.True(t, err.IsOk())

	ctx.Set(testResource.ProviderIDKey, provider.IdentifierID)

	return provider
}

func NetworkFixture(ctx context.Context, t *testing.T, nuc usecase.Network) *network.Network {
	parentID := ctx.Value(testResource.ProviderIDKey).(identifier.Provider)
	netName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))
	req := dto.CreateNetworkRequest{
		ParentIDProvider: &parentID,
		Name:             netName,
		Managed:          false,
	}
	ctx.Set(context.RequestKey, req)

	network := &network.Network{}
	err := nuc.Create(ctx, network)

	ctx.Set(testResource.NetworkIDKey, network.IdentifierID)
	require.Equal(t, errors.Created.Code, err.Code)

	return network
}

func SubnetworkFixture(ctx context.Context, t *testing.T, suc usecase.Subnetwork) *network.Subnetwork {
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

	subnetwork := &network.Subnetwork{}
	err := suc.Create(ctx, subnetwork)
	ctx.Set(testResource.SubnetworkIDKey, subnetwork.IdentifierID)
	require.Equal(t, errors.Created.Code, err.Code)

	return subnetwork
}

func VMFixture(ctx context.Context, t *testing.T, vuc usecase.VM) *instance.VM {
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
		OS: instance.VMOS{
			ID: osName,
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
	require.Equal(t, errors.Created, err)

	return vm
}

func SqlDBFixture(ctx context.Context, t *testing.T, vuc usecase.SqlDB) *instance.SqlDB {
	region := "europe-west9"

	req := dto.CreateSqlDBRequest{
		ParentID:    ctx.Value(testResource.NetworkIDKey).(identifier.Network),
		Name:        "test-sqldb",
		Region:      region,
		MachineType: "db-f1-micro",
		Disk: instance.SqlDbDisk{
			Type:    instance.DiskTypeSSD,
			SizeGib: 10,
		},
		SQLType:    instance.PostgresSQLType,
		SQLVersion: "12",

		Managed: true,
	}
	ctx.Set(context.RequestKey, req)

	db := &instance.SqlDB{}
	err := vuc.Create(ctx, db)
	require.Equal(t, errors.Created, err)

	return db
}

func FirewallFixture(ctx context.Context, t *testing.T, fuc usecase.Firewall) *network.Firewall {
	parentID := ctx.Value(testResource.NetworkIDKey).(identifier.Network)
	firewallName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))

	req := dto.CreateFirewallRequest{
		ParentID: parentID,
		Name:     firewallName,
		AllowRules: network.FirewallRuleList{
			{
				Protocol: "tcp",
				Ports:    []string{"80", "443"},
			},
			//{
			//	Protocol: "udp",
			//	Ports:    []string{"53"},
			//},
		},
		DenyRules: network.FirewallRuleList{
			{
				Protocol: "tcp",
				Ports:    []string{"65"},
			},
		},
		Managed: false,
	}
	ctx.Set(context.RequestKey, req)
	firewall := &network.Firewall{}
	err := fuc.Create(ctx, firewall)
	require.Equal(t, errors.Created.Code, err.Code)
	return firewall
}
