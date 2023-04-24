package plugin

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
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
//if err1 := plugin.ValidateAndCompletePluginEntry(pluginEntry); !err1.IsOk() {
//logger.Info.Println(err1)
//}

func TestValidateWithSubtypesShouldReturnOK(t *testing.T) {
	pluginEntry := map[string]interface{}{
		"vpc":         "vpc-1234567890",
		"zone":        "us-east-1a",
		"network":     "network-1234567890",
		"subnetwork":  "subnet-1234567890",
		"machineType": "n1-standard-1",
		"disk": map[string]interface{}{
			"size":       10,
			"type":       "pd-standard",
			"mode":       "READ_WRITE",
			"autoDelete": true,
		},
		"os": map[string]interface{}{
			"type":    "ubuntu",
			"version": "20.04",
		},
	}
	plugin, err := Get("gcp", "vm")
	if !err.IsOk() {
		t.Errorf("Error getting plugin: %v", err)
	}
	//logger.Info.Println(plugin.Types[0])
	if _, err1 := plugin.ValidateAndCompletePluginEntry(pluginEntry); !err1.IsOk() {
		t.Errorf("Error validating plugin: %v", err1)
	}
}

func TestValidationWithSubtypeShouldReturnMissingValue(t *testing.T) {
	pluginEntry := map[string]interface{}{
		"vpc":         "vpc-1234567890",
		"zone":        "us-east-1a",
		"network":     "network-1234567890",
		"subnetwork":  "subnet-1234567890",
		"machineType": "n1-standard-1",
		"disk": map[string]interface{}{
			"size":       10,
			"type":       "pd-standard",
			"mode":       "READ_WRITE",
			"autoDelete": true,
		},
		"os": map[string]interface{}{
			"type": "ubuntu",
		},
	}
	plugin, err := Get("gcp", "vm")
	if !err.IsOk() {
		t.Errorf("Error getting plugin: %v", err)
	}
	//logger.Info.Println(plugin.Types[0])
	if _, err1 := plugin.ValidateAndCompletePluginEntry(pluginEntry); err1.IsOk() {
		logger.Info.Println(err1)
		t.Errorf("Validation should have failed")
	}
}
