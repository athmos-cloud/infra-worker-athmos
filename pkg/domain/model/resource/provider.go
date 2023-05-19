package resource

import (
	"fmt"
	identifier2 "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
)

const (
	providerIDSuffixLength = 10
)

type Provider struct {
	Metadata   metadata.Metadata    `bson:"metadata"`
	Identifier identifier2.Provider `bson:"identifier"`
	Auth       secret.Secret        `bson:"secret" plugin:"secret"`
	VPCs       VPCCollection        `bson:"vpcs"`
	Networks   NetworkCollection    `bson:"networks"`
}

type ProviderCollection map[string]Provider

func NewProvider(payload NewResourcePayload) Provider {
	payload.Validate()
	id := identifier2.Provider{
		ProviderID: fmt.Sprintf("%s-%s", formatResourceName(payload.Name), utils.RandomString(providerIDSuffixLength)),
	}
	return Provider{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name:         id.ProviderID,
			NotMonitored: !payload.Managed,
			Tags:         payload.Tags,
		}),
		Identifier: id,
		VPCs:       make(VPCCollection),
		Networks:   make(NetworkCollection),
	}
}
