package kubernetes

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/share"
	"k8s.io/client-go/dynamic"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

var kubeClient *Kubernetes
var lock = &sync.Mutex{}

type Kubernetes struct {
	DynamicClient dynamic.Interface
	K8sClient     client.Client
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	if kubeClient == nil && os.Getenv(share.EnvTypeEnvironmentVariable) != string(share.EnvTypeTest) {
		logger.Info.Printf("Init kubernetes client...")
		conf := ctrl.GetConfigOrDie()
		dynamicCli := dynamic.NewForConfigOrDie(conf)

		cfg, err := ctrl.GetConfig()
		if err != nil {
			panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error getting Kubernetes config: %v", err)))
		}
		k8sClient, err := client.New(cfg, client.Options{})
		if err != nil {
			panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error getting Kubernetes client: %v", err)))
		}

		SetClient(&Kubernetes{
			DynamicClient: dynamicCli,
			K8sClient:     k8sClient,
		})
	}
}
func Client() *Kubernetes {
	return kubeClient
}
func SetClient(cli *Kubernetes) {
	kubeClient = cli
}
