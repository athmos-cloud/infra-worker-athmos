package plugin

import (
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/plugin"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	plugin2 "github.com/PaulBarrie/infra-worker/pkg/plugin"
)

type Service struct{}

func (service *Service) GetPlugin(request dto.GetPluginRequest) (dto.GetPluginResponse, errors.Error) {
	plugin, err := plugin2.Get(request.ProviderType, request.ResourceType)
	if !err.IsOk() {
		return dto.GetPluginResponse{}, err
	}
	return dto.GetPluginResponse{Plugin: plugin}, errors.OK
}
