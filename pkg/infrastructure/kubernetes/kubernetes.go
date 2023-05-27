package kubernetes

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/share"
	"os"
	providerGCP "github.com/upbound/provider-gcp/apis/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/dynamic"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	"sync"
)

var kubeClient client.Client
var lock = &sync.Mutex{}

func init() {
	lock.Lock()
	defer lock.Unlock()
	if kubeClient == nil && os.Getenv(share.EnvTypeEnvironmentVariable) != string(share.EnvTypeTest) {
		logger.Info.Printf("Init kubernetes client...")
		cliSchema := getSchema()
		cfg, err := ctrl.GetConfig()
		if err != nil {
			panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error getting Kubernetes config: %v", err)))
		}
		k8sClient, err := client.New(cfg, client.Options{Scheme: cliSchema})
		if err != nil {
			panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error getting Kubernetes client: %v", err)))
		}
		SetClient(k8sClient)
	}
}
func Client() client.Client {
	return kubeClient
}
func SetClient(cli client.Client) {
	kubeClient = cli
}

func getSchema() *runtime.Scheme {
	newScheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(newScheme))
	registerGCPResources(newScheme)

	return newScheme
}

func registerGCPResources(runtimeScheme *runtime.Scheme) {
	SchemeBuilder := &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "gcp.upbound.io", Version: "v1beta1"}}
	SchemeBuilder.Register(&providerGCP.ProviderConfig{}, &providerGCP.ProviderConfigList{})
	err := SchemeBuilder.AddToScheme(runtimeScheme)
	if err != nil {
		panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error registering GCP resources: %v", err)))
	}
}
