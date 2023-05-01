package kubernetes

import (
	"bufio"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

const (
	manifestSeparatorString        = "---"
	groupAPIVersionSeparatorString = "/"
)

type Identifier struct {
	ResourceID schema.GroupVersionResource
	Name       string
	Namespace  string
}

func (identifier Identifier) Equals(other Identifier) bool {
	return identifier.ResourceID == other.ResourceID &&
		identifier.Name == other.Name &&
		identifier.Namespace == other.Namespace
}

func GetResourcesIdentifiersFromManifests(manifests string) []Identifier {
	var identifierList []Identifier
	manifestList := strings.Split(manifests, manifestSeparatorString)
	for _, val := range manifestList {
		if val == "" {
			continue
		}
		identifier := getResourceIdentifierFromManifest(val)
		identifierList = append(identifierList, identifier)
	}
	return identifierList
}

func getResourceIdentifierFromManifest(manifest string) Identifier {
	formatManifestString(&manifest)
	var out map[string]interface{}
	if err := yaml.Unmarshal([]byte(manifest), &out); err != nil {
		panic(errors.InternalError.WithMessage(err.Error()))
	}
	var apiVersion string
	var group string
	var kind string
	var name string
	var namespace string

	if val, ok := out["kind"]; ok {
		kind = val.(string)
	}
	if val, ok := out["apiVersion"]; ok {
		groupAPI := strings.Split(val.(string), groupAPIVersionSeparatorString)
		group = groupAPI[0]
		apiVersion = groupAPI[1]
	}
	if val, ok := out["metadata"].(map[string]interface{})["name"]; ok {
		name = val.(string)
	}
	if val, ok := out["metadata"].(map[string]interface{})["namespace"]; ok {
		namespace = val.(string)
	}
	if group == "" || apiVersion == "" || kind == "" || name == "" {
		panic(errors.InvalidArgument.WithMessage("Manifest is not valid"))
	}
	return Identifier{
		ResourceID: schema.GroupVersionResource{
			Group:    group,
			Version:  apiVersion,
			Resource: kind,
		},
		Name:      name,
		Namespace: namespace,
	}
}

func formatManifestString(manifest *string) {
	newManifest := strings.ReplaceAll(*manifest, manifestSeparatorString, "")
	var output strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(newManifest))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "#") {
			continue
		}
		output.WriteString(line)
		output.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
	*manifest = output.String()
}
