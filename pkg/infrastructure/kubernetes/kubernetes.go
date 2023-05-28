package kubernetes

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/share"
	"k8s.io/client-go/dynamic"
	"log"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

var kubeClient *Kubernetes
var lock = &sync.Mutex{}

type Kubernetes struct {
	Client        client.Client
	DynamicClient *dynamic.DynamicClient
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	if kubeClient == nil && os.Getenv(share.EnvTypeEnvironmentVariable) != string(share.EnvTypeTest) {
		logger.Info.Printf("Init kubernetes client...")
		cliSchema := getScheme()
		cfg, err := ctrl.GetConfig()
		if err != nil {
			panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error getting Kubernetes config: %v", err)))
		}
		k8sClient, err := client.New(cfg, client.Options{Scheme: cliSchema})
		if err != nil {
			panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error getting Kubernetes client: %v", err)))
		}
		dynamicClient, err := dynamic.NewForConfig(cfg)
		if err != nil {
			log.Fatal(err)
		}
		SetClient(&Kubernetes{Client: k8sClient, DynamicClient: dynamicClient})
	}
}
func Client() *Kubernetes {
	return kubeClient
}
func SetClient(cli *Kubernetes) {
	kubeClient = cli
}
