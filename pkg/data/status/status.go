package status

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/helm"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
)

const randomNameUUIDLength = 5

type ResourceStatus struct {
	HelmRelease         helm.ReleaseReference   `bson:"helmRelease"`
	KubernetesResources kubernetes.ResourceList `bson:"kubernetesResources"`
	PluginReference     plugin.Reference        `bson:"pluginReference"`
}

func (status *ResourceStatus) Equals(other ResourceStatus) bool {
	return status.HelmRelease.Equals(other.HelmRelease) &&
		status.KubernetesResources.Equals(other.KubernetesResources)
}

func New(name string, resourceType types.ResourceType, providerType types.ProviderType) ResourceStatus {
	ref := plugin.NewReference(resourceType, providerType)
	return ResourceStatus{
		HelmRelease:         helm.NewRelease(fmt.Sprintf("%s-%s", name, utils.RandomString(randomNameUUIDLength))),
		KubernetesResources: make(kubernetes.ResourceList, 0),
		PluginReference:     ref,
	}
}
