package application

import (
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/plugin"
	plugin2 "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type PluginService struct{}

func (service *PluginService) GetPlugin(request dto.GetPluginRequest) (dto.GetPluginResponse, errors.Error) {
	plugin, err := plugin2.Get(request.ProviderType, request.ResourceType)
	if !err.IsOk() {
		return dto.GetPluginResponse{}, err
	}
	return dto.GetPluginResponse{Plugin: plugin}, errors.OK
}
