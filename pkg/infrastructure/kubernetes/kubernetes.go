package kubernetes

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/share"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

var kubeClient client.Client
var lock = &sync.Mutex{}

func init() {
	lock.Lock()
	defer lock.Unlock()
	if kubeClient == nil && os.Getenv(share.EnvTypeEnvironmentVariable) != string(share.EnvTypeTest) {
		logger.Info.Printf("Init kubernetes client...")
		cfg, err := ctrl.GetConfig()
		if err != nil {
			panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error getting Kubernetes config: %v", err)))
		}
		k8sClient, err := client.New(cfg, client.Options{})
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
