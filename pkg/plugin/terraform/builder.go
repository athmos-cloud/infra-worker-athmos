package terraform

import (
	"context"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository"
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository/bucket"
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/types/workdir"
	"github.com/PaulBarrie/infra-worker/pkg/plugin"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/common"
	tfConf "github.com/PaulBarrie/infra-worker/pkg/plugin/terraform/config/module"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/terraform/config/variable"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/terraform/config/provider"
	"reflect"
)

const (
	ModuleFolderName       = "modules"
	MainTerraformFile      = "main.tf"
	ProvidersTerraformFile = "providers.tf"
	VariablesTerraformFile = "variables.tf"
	VariableInputFiles     = "variables.tfvars"
	OutputsTerraformFile   = "outputs.tf"
)

type Builder struct {
	PluginLibraryRepository repository.IRepository
	DependencyRepository    repository.IRepository
	CurrentStrategy         plugin.IPluginStrategy
}

func NewBuilder(pluginLibraryRepository repository.IRepository, dependencyRepository repository.IRepository) *Builder {
	return &Builder{
		PluginLibraryRepository: pluginLibraryRepository,
		DependencyRepository:    dependencyRepository,
	}
}

func (pb *Builder) Load(ctx context.Context, optn option.Option) (common.Plugin, errors.Error) {
	if !optn.SetType(reflect.String).Validate() {
		return common.Plugin{}, errors.New("invalid option type")
	}
	id := optn.Value.(string)
	mongoRes, err := pb.PluginLibraryRepository.Get(
		ctx,
		option.Option{
			Type:  reflect.TypeOf(reflect.TypeOf(mongo.RetrieveRequestPayload{})).Kind(),
			Value: mongo.RetrieveRequestPayload{Id: id},
		})
	logger.Info.Printf("plugin_: %v", mongoRes)
	plugin_ := mongoRes.(common.Plugin)
	switch plugin_.Type {
	case common.Terraform:
		pb.CurrentStrategy = NewStrategy(plugin_)
	}
	return plugin_, err
}

func (pb *Builder) Vendor(ctx context.Context, plugin common.Plugin) (workdir.Workdir, errors.Error) {
	module, err := pb.DependencyRepository.Get(
		ctx,
		option.Option{
			Type: reflect.TypeOf(bucket.RetrieveRequestPayload{}).Kind(),
			Value: bucket.RetrieveRequestPayload{
				BucketName: plugin.Location.Bucket,
				Dir:        plugin.Location.Folder,
			},
		})
	if !err.IsOk() {
		return workdir.Workdir{}, err
	}
	moduleFolder := module.(workdir.Folder)
	return workdir.Workdir{
		Folders: []workdir.Folder{
			{
				Name: ModuleFolderName,
				Folders: []workdir.Folder{
					{
						Name:    plugin.Name,
						Folders: moduleFolder.Folders,
						Files:   moduleFolder.Files,
					},
				},
			},
		},
	}, errors.OK
}

func (pb *Builder) Build(ctx context.Context, payload plugin.BuildPluginPayload) (workdir.Workdir, errors.Error) {
	// Parse module content
	plugin := payload.PluginCall.Plugin
	// Get variable payload
	terraformVariables := variable.ListFromCommon(plugin.Config.Inputs)
	packagedVariables := terraformVariables.Build(payload.PluginCall.Name, payload.PluginCall.Inputs)
	// Write the module call -> main.tf
	moduleRefList := variable.ReferenceList{}
	for _, v := range *packagedVariables {
		moduleRefList = append(moduleRefList, *v.ModuleReference)
	}
	module := tfConf.Module{
		Name:      tfConf.GetModuleName(plugin.Name, payload.PluginCall.Name),
		Source:    fmt.Sprintf("./%s/%s", ModuleFolderName, plugin.Name),
		Variables: moduleRefList,
	}
	payload.Workdir.Files = append(payload.Workdir.Files, workdir.File{
		Name:    MainTerraformFile,
		Content: module.ToString(),
	})
	// Set the providers -> providers.tf
	provider := provider.FromCommon(plugin.Config.Dependencies)
	payload.Workdir.Files = append(payload.Workdir.Files, workdir.File{
		Name:    ProvidersTerraformFile,
		Content: ,
	}
	// Set the outputs -> outputs.tf

}

func (pb *Builder) Register(ctx context.Context, location common.Location) (common.Plugin, errors.Error) {
	//TODO implement me
	panic("implement me")
}
