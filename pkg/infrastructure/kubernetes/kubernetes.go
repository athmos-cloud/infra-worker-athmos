package kubernetes

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

var Client *Kubernetes
var lock = &sync.Mutex{}

type Kubernetes struct {
	DynamicClient dynamic.Interface
	K8sClient     client.Client
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	if Client == nil {
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

		Client = &Kubernetes{
			DynamicClient: dynamicCli,
			K8sClient:     k8sClient,
		}
	}
}
