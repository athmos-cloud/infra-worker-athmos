package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws/xrds"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	instanceModel "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	networkModel "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	awsCompute "github.com/upbound/provider-aws/apis/ec2/v1beta1"
	awsNetworks "github.com/upbound/provider-aws/apis/networkfirewall/v1beta1"
	"github.com/upbound/provider-aws/apis/v1beta1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"
)

func ProviderFixture(ctx context.Context, t *testing.T, puc usecase.Provider) *resource.Provider {
	req := dto.CreateProviderRequest{
		Name:           "fixture-provider",
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

func NetworkFixture(ctx context.Context, t *testing.T, nuc usecase.Network) *networkModel.Network {
	parentID := ctx.Value(testResource.ProviderIDKey).(identifier.Provider)
	//netName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))
	req := dto.CreateNetworkRequest{
		ParentIDProvider: &parentID,
		Name:             "fixture-network",
		Region:           "eu-west-1",
		Managed:          false,
	}
	ctx.Set(context.RequestKey, req)

	network := &networkModel.Network{}
	err := nuc.Create(ctx, network)

	ctx.Set(testResource.NetworkIDKey, network.IdentifierID)
	require.Equal(t, errors.Created.Code, err.Code)

	return network
}

func SubnetworkFixture(ctx context.Context, t *testing.T, suc usecase.Subnetwork) *networkModel.Subnetwork {
	parentID := ctx.Value(testResource.NetworkIDKey).(identifier.Network)
	//subnetName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))
	region := "eu-west-1"
	ipCIDR := "10.0.0.1/26"
	req := dto.CreateSubnetworkRequest{
		ParentID:    parentID,
		Name:        "fixture-subnet",
		Region:      region,
		IPCIDRRange: ipCIDR,
		Managed:     false,
	}
	ctx.Set(context.RequestKey, req)

	subnetwork := &networkModel.Subnetwork{}
	err := suc.Create(ctx, subnetwork)
	ctx.Set(testResource.SubnetworkIDKey, subnetwork.IdentifierID)
	require.Equal(t, errors.Created.Code, err.Code)

	return subnetwork
}

func VMFixture(ctx context.Context, t *testing.T, vuc usecase.VM) *instanceModel.VM {
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
		OS: instanceModel.VMOS{
			ID:   osName,
			Name: osName,
		},
		Disks: []instanceModel.VMDisk{
			{
				AutoDelete: autoDelete,
				Mode:       instanceModel.DiskModeReadWrite,
				Type:       instanceModel.DiskTypeSSD,
				SizeGib:    10,
			},
		},
		Managed: true,
	}
	ctx.Set(context.RequestKey, req)

	vm := &instanceModel.VM{}
	err := vuc.Create(ctx, vm)

	require.Equal(t, errors.Created, err)

	return vm
}

func FirewallFixture(ctx context.Context, t *testing.T, fuc usecase.Firewall) *networkModel.Firewall {
	parentID := ctx.Value(testResource.NetworkIDKey).(identifier.Network)
	req := dto.CreateFirewallRequest{
		ParentID: parentID,
		Name:     "fixture-firewall",
		AllowRules: networkModel.FirewallRuleList{
			{
				Protocol: "tcp",
				Ports:    []string{"80", "443"},
			},
			{
				Protocol: "udp",
				Ports:    []string{"53"},
			},
		},
		DenyRules: networkModel.FirewallRuleList{
			{
				Protocol: "tcp",
				Ports:    []string{"65"},
			},
		},
		Managed: false,
	}
	ctx.Set(context.RequestKey, req)
	firewall := &networkModel.Firewall{}
	err := fuc.Create(ctx, firewall)

	require.Equal(t, errors.Created.Code, err.Code)
	return firewall
}

func SqlDBFixture(ctx context.Context, t *testing.T, vuc usecase.SqlDB) *instanceModel.SqlDB {
	region := "eu-west-1"
	req := dto.CreateSqlDBRequest{
		ParentID:    ctx.Value(testResource.NetworkIDKey).(identifier.Network),
		Name:        "fixture-db",
		Region:      region,
		MachineType: "db.m7g",
		Disk: instanceModel.SqlDbDisk{
			Type:    instanceModel.DiskTypeSSD,
			SizeGib: 10,
		},
		SQLType:      instanceModel.PostgresSQLType,
		SQLVersion:   "12",
		Managed:      true,
		RootPassword: "proEsgi7656$!",
	}
	ctx.Set(context.RequestKey, req)

	db := &instanceModel.SqlDB{}
	err := vuc.Create(ctx, db)
	require.Equal(t, errors.Created, err)

	return db
}

func ClearFixtures(ctx context.Context) {
	ClearSqlFixtures(ctx)
	ClearVMFixtures(ctx)
	ClearSubnetworkFixtures(ctx)
	ClearFirewallFixtures(ctx)
	ClearNetworksFixtures(ctx)
	ClearProviderFixtures(ctx)
}

func ClearSqlFixtures(ctx context.Context) {
	namespace := ctx.Value(test.TestNamespaceContextKey).(string)

	rdsInstances := &xrds.SQLDatabaseList{}
	err := kubernetes.Client().Client.List(ctx, rdsInstances, &client.ListOptions{
		Namespace: namespace,
	})

	if err != nil {
		return
	}
	for _, rdsInstance := range rdsInstances.Items {
		err = kubernetes.Client().Client.Delete(ctx, &rdsInstance)
		if err != nil {
			logger.Warning.Println(fmt.Sprintf("Error deleting sql db %s", rdsInstance.Name))
			continue
		}
	}

	passwordSecrets := &v1.SecretList{}
	err = kubernetes.Client().Client.List(ctx, passwordSecrets, &client.ListOptions{
		Namespace: namespace,
	})
	if err != nil {
		return
	}
	for _, passwordSecret := range passwordSecrets.Items {
		err = kubernetes.Client().Client.Delete(ctx, &passwordSecret)
		if err != nil {
			logger.Warning.Println(fmt.Sprintf("Error deleting password secret %s", passwordSecret.Name))
			continue
		}
	}

	logger.Info.Println("Cleared Sql fixtures")
	time.Sleep(5 * 1000 * time.Millisecond)
}

func ClearVMFixtures(ctx context.Context) {
	keyPairs := &awsCompute.KeyPairList{}
	err := kubernetes.Client().Client.List(ctx, keyPairs)
	if err != nil {
		return
	}
	for _, keyPair := range keyPairs.Items {
		err = kubernetes.Client().Client.Delete(ctx, &keyPair)
		if err != nil {
			logger.Warning.Printf("Error deleting key pair %s: %v", keyPair.Name, err)
			continue
		}
	}

	ec2Instances := &awsCompute.InstanceList{}
	err = kubernetes.Client().Client.List(ctx, ec2Instances)
	if err != nil {
		return
	}
	for _, ec2Instance := range ec2Instances.Items {
		err = kubernetes.Client().Client.Delete(ctx, &ec2Instance)
		if err != nil {
			logger.Warning.Printf("Error deleting ec2 instance %s: %v", ec2Instance.Name, err)
			continue
		}
	}
}

func ClearSubnetworkFixtures(ctx context.Context) {
	subnetworks := &awsCompute.SubnetList{}
	err := kubernetes.Client().Client.List(ctx, subnetworks)
	if err != nil {
		return
	}
	for _, subnetwork := range subnetworks.Items {
		err = kubernetes.Client().Client.Delete(ctx, &subnetwork)
		if err != nil {
			logger.Warning.Printf("Error deleting subnetwork %s: %v", subnetwork.Name, err)
			continue
		}
	}
}

func ClearFirewallFixtures(ctx context.Context) {
	ruleGroups := &awsNetworks.RuleGroupList{}
	err := kubernetes.Client().Client.List(ctx, ruleGroups)
	if err != nil {
		return
	}
	for _, ruleGroup := range ruleGroups.Items {
		err = kubernetes.Client().Client.Delete(ctx, &ruleGroup)
		if err != nil {
			logger.Warning.Printf("Error deleting firewall %s: %v", ruleGroup.Name, err)
			continue
		}
	}

	firewallPolicies := &awsNetworks.FirewallPolicyList{}
	err = kubernetes.Client().Client.List(ctx, firewallPolicies)
	if err != nil {
		return
	}
	for _, firewallPolicy := range firewallPolicies.Items {
		err = kubernetes.Client().Client.Delete(ctx, &firewallPolicy)
		if err != nil {
			logger.Warning.Printf("Error deleting firewall %s: %v", firewallPolicy.Name, err)
			continue
		}
	}

	firewalls := &awsNetworks.FirewallList{}
	err = kubernetes.Client().Client.List(ctx, firewalls)
	if err != nil {
		return
	}
	for _, firewall := range firewalls.Items {
		err = kubernetes.Client().Client.Delete(ctx, &firewall)
		if err != nil {
			logger.Warning.Printf("Error deleting firewall %s: %v", firewall.Name, err)
			continue
		}
	}
}

func ClearNetworksFixtures(ctx context.Context) {
	networks := &awsCompute.VPCList{}
	err := kubernetes.Client().Client.List(ctx, networks)
	if err != nil {
		return
	}
	for _, network := range networks.Items {
		err = kubernetes.Client().Client.Delete(ctx, &network)
		if err != nil {
			logger.Warning.Printf("Error deleting network %s: %v", network.Name, err)
			continue
		}
	}
}

func ClearProviderFixtures(ctx context.Context) {
	providers := &v1beta1.ProviderConfigList{}
	err := kubernetes.Client().Client.List(ctx, providers)
	if err != nil {
		return
	}
	for _, provider := range providers.Items {
		err = kubernetes.Client().Client.Delete(ctx, &provider)
		if err != nil {
			logger.Warning.Printf("Error deleting provider %s: %v", provider.Name, err)
			continue
		}
	}
}
