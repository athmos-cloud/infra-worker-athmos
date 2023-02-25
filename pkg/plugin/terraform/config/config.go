package config

import (
	"bytes"
	"encoding/gob"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/utils"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/common"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/terraform/config/module"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/terraform/config/output"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/terraform/config/provider"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/terraform/config/variable"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"os"
)

type Config struct {
	Modules   module.ModuleList `hcl:"module,block" bson:"modules"`
	Variables variable.List     `hcl:"variable,block" bson:"variable"`
	Providers provider.List     `hcl:"provider,block" bson:"providers"`
	Output    output.List       `bson:"provider,output"`
}

type StringPayload struct {
	Modules   string `bson:"modules"`
	Variables string `bson:"variable"`
	Providers string `bson:"providers"`
}

type List []Config

func FromCommon(plugin *common.Plugin) *Config {
	return &Config{}
}

func Marshall(payload string) (*Config, errors.Error) {
	var res Config
	file, err := utils.StringToTempFile(payload)
	if !err.IsOk() {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	errDec := hclsimple.DecodeFile(file.Name(), nil, res)
	if errDec != nil {
		return nil, errors.ParseError.WithMessage(err)
	}
	return &res, errors.OK
}

func (c *Config) Hash() string {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(c)
	if err != nil {
		logger.Warning.Printf("Error encoding output %s", err)
		return ""
	}
	return string(b.Bytes())
}

func (c *Config) Merge(config Config, optn option.Option) {
	c.Modules.MergeList(config.Modules, optn)
	c.Variables.MergeList(config.Variables, optn)
	c.Providers.MergeList(config.Providers, optn)
	c.Output.MergeList(config.Output, optn)
}

func (c *Config) ToString() StringPayload {
	return StringPayload{
		Modules:   c.Modules.ToString(),
		Variables: c.Variables.ToString(),
		Providers: c.Providers.ToString(),
	}
}

//func (c *Config) ToString() (string, error) {
//	if err := c.Generate(); err != nil {
//		logger.Error.Println("Error generating module: ", err)
//		return "", err
//	}
//	return c.PluginConfig.CodeContent, nil
//}
//
//func (c *Config) Get() *config.PluginConfig {
//	return c.PluginConfig
//}
//
//func (c *Config) Generate() error {
//	buffer := new(bytes.Buffer)
//	templateFile, err := template.ParseFiles(terraform.ModuleTemplateFileLocation)
//	if err != nil {
//		logger.Error.Println("parse template_file: ", err)
//		return err
//	}
//
//	configMap := map[string]interface{}{
//		"name":   c.PluginConfig.Name,
//		"source": c.PluginConfig.ModulePath,
//		"values": c.PluginConfig.ModuleVars,
//	}
//	for _, dependency := range *c.Vendor {
//		err, _ := dependency.Load()
//		if err != nil {
//			return err
//		}
//	}
//
//	if err = templateFile.Execute(buffer, configMap); err != nil {
//		logger.Error.Println("execute template_file: ", err)
//		return err
//	}
//	c.PluginConfig.CodeContent = buffer.String()
//	return nil
//}
