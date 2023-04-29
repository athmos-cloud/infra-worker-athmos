package application

import (
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/plugin"
	plugin2 "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type PluginService struct{}

func (service *PluginService) GetPlugin(request dto.GetPluginRequest) (dto.GetPluginResponse, errors.Error) {
	err := errors.OK
	defer func() {
		if r := recover(); r != nil {
			err = r.(errors.Error)
		}
	}()
	plugin := plugin2.Get(plugin2.ResourceReference{ResourceType: request.ResourceType, ProviderType: request.ProviderType})

	return dto.GetPluginResponse{Plugin: plugin}, err
}
