package plugin

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository/bucket"
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/types/workdir"
	common "github.com/PaulBarrie/infra-worker/pkg/plugin/common"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/terraform"
)

type IPluginBuilder interface {
	Load(context.Context, option.Option) (common.Plugin, errors.Error)
	Vendor(context.Context, common.Plugin) (workdir.Workdir, errors.Error)
	Build(context.Context, BuildPluginPayload) (workdir.Workdir, errors.Error)
	Register(context.Context, common.Location) (common.Plugin, errors.Error)
}

type BuildPluginPayload struct {
	Workdir    workdir.Workdir
	PluginCall common.PluginInstance
}

func FactoryBuilder(pluginType common.Type) IPluginBuilder {
	switch pluginType {
	case common.Terraform:
		return terraform.NewBuilder(mongo.Client, bucket.MinioClient)
	}
	return nil
}
