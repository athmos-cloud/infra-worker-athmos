package output

import (
	"bytes"
	"encoding/gob"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	common "github.com/PaulBarrie/infra-worker/pkg/plugin/common/config"
	"github.com/hashicorp/hcl/v2"
	"reflect"
)

type Output struct {
	Name        string   `hcl:",label"`
	Description string   `hcl:"description,optional"`
	Sensitive   bool     `hcl:"sensitive,optional"`
	Value       string   `hcl:"value,optional"`
	Options     hcl.Body `hcl:",remain"`
}

func (o *Output) Hash() string {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(o)
	if err != nil {
		logger.Warning.Printf("Error encoding output %s", err)
		return ""
	}
	return string(b.Bytes())
}

func (o *Output) ToCommon(value interface{}) *common.Output {
	return &common.Output{
		Name:  o.Name,
		Value: value,
		Type:  reflect.TypeOf(value).Kind().String(),
	}
}
func (o *Output) ToString() string {
	return `
		output "` + o.Name + `" {
			value = ` + o.Value + `
		}
	`
}
