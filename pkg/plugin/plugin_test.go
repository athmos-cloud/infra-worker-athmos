package plugin

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"reflect"
	"testing"
)

//pluginEntry := map[string]interface{}{
//"vpc":         "vpc-1234567890",
//"zone":        "us-east-1a",
//"network":     "network-1234567890",
//"subnetwork":  "subnet-1234567890",
//"machineType": "n1-standard-1",
//"disk": map[string]interface{}{
//"size":       10,
//"type":       "pd-standard",
//"mode":       "READ_WRITE",
//"autoDelete": true,
//},
//"os": map[string]interface{}{
//"type":    "ubuntu",
//"version": "20.04",
//},
//}
//plugin, err := plugin.Get("gcp", "vm")
//if !err.IsOk() {
//panic(err)
//}
////logger.Info.Println(plugin.Types[0])
//if err1 := plugin.Validate(pluginEntry); !err1.IsOk() {
//logger.Info.Println(err1)
//}

func TestGet(t *testing.T) {
	type args struct {
		provider     common.ProviderType
		resourceType common.ResourceType
	}
	tests := []struct {
		name  string
		args  args
		want  Plugin
		want1 errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Get(tt.args.provider, tt.args.resourceType)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInput_Validate(t *testing.T) {
	type fields struct {
		Name        string
		Description string
		Type        string
		Default     interface{}
		Required    bool
	}
	type args struct {
		entry map[string]interface{}
		types []Type
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Input{
				Name:        tt.fields.Name,
				Description: tt.fields.Description,
				Type:        tt.fields.Type,
				Default:     tt.fields.Default,
				Required:    tt.fields.Required,
			}
			if got := i.Validate(tt.args.entry, tt.args.types); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlugin_Validate(t *testing.T) {
	type fields struct {
		Prerequisites []Prerequisite
		Inputs        []Input
		Types         []Type
	}
	type args struct {
		entry map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{
				Prerequisites: tt.fields.Prerequisites,
				Inputs:        tt.fields.Inputs,
				Types:         tt.fields.Types,
			}
			if got := p.Validate(tt.args.entry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
