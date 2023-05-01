package kubernetes

import (
	"context"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sync"
)

const (
	kubeConfigEnvVar = "KUBECONFIG"
)

var Client *DAO
var lock = &sync.Mutex{}

type DAO struct {
	DynamicClient dynamic.Interface
	ClientSet     *kubernetes.Clientset
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	if Client == nil {
		config := ctrl.GetConfigOrDie()
		dynamicCli := dynamic.NewForConfigOrDie(config)

		var restConfig *rest.Config
		var err error
		if kubeConfig := os.Getenv(kubeConfigEnvVar); kubeConfig != "" {
			restConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		} else {
			restConfig, err = rest.InClusterConfig()
		}
		if err != nil {
			panic(fmt.Sprintf("Error getting Kubernetes configuration: %v\n", err))
		}
		clientSet, err := kubernetes.NewForConfig(restConfig)
		if err != nil {
			fmt.Printf("Error creating Kubernetes clientSet: %v\n", err)
			os.Exit(1)
		}
		Client = &DAO{
			DynamicClient: dynamicCli,
			ClientSet:     clientSet,
		}
	}
}

func (dao *DAO) Get(ctx context.Context, option option.Option) interface{} {
	if option = option.SetType(reflect.TypeOf(GetResourceRequest{}).String()); !option.Validate() {
		panic(errors.InvalidArgument.WithMessage(
			"Argument must be a kubernetes.GetResourceRequest{resourceId, namespace, name}",
		))
	}
	request := option.Get().(GetResourceRequest)

	resource, err := dao.DynamicClient.Resource(request.ResourceID).
		Namespace(request.Namespace).Get(ctx, request.Name, metav1.GetOptions{})
	if err != nil {
		panic(errors.NotFound.WithMessage(err.Error()))
	}

	return resource
}

func (dao *DAO) Exists(ctx context.Context, o option.Option) bool {
	//TODO implement me
	panic("implement me")
}

func (dao *DAO) GetAll(ctx context.Context, option option.Option) interface{} {
	if option = option.SetType(reflect.TypeOf(GetListResourceRequest{}).String()); !option.Validate() {
		panic(errors.InvalidArgument.WithMessage(
			"Argument must be a kubernetes.GetListResourceRequest{resourceId, namespace, labels}",
		))
	}
	request := option.Get().(GetListResourceRequest)
	var options metav1.ListOptions
	if request.Labels != nil {
		options = metav1.ListOptions{
			LabelSelector: labelsToString(request.Labels),
		}
	}
	list, err := dao.DynamicClient.Resource(request.ResourceID).Namespace(request.Namespace).List(ctx, options)
	if err != nil {
		panic(errors.NotFound.WithMessage(err.Error()))
	}
	return list.Items
}

// Create namespace from a given name string
func (dao *DAO) Create(ctx context.Context, opt option.Option) interface{} {
	if opt = opt.SetType(reflect.TypeOf(CreateNamespaceRequest{}).String()); opt.Validate() {
		return dao.createNamespace(ctx, opt.Get().(CreateNamespaceRequest).Name)
	} else if opt = opt.SetType(reflect.TypeOf(CreateSecretRequest{}).String()); opt.Validate() {
		return dao.createSecret(ctx, opt.Get().(CreateSecretRequest))
	} else {
		panic(errors.InvalidArgument.WithMessage(
			"Argument must be a kubernetes.CreateNamespaceRequest{name}",
		))
	}

}

func (dao *DAO) createNamespace(ctx context.Context, namespace string) interface{} {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	createdNamespace, err := dao.ClientSet.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil {
		panic(errors.ExternalServiceError.WithMessage(err.Error()))
	}
	return createdNamespace

}

func (dao *DAO) createSecret(ctx context.Context, request CreateSecretRequest) interface{} {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: request.Name,
		},
		Data: map[string][]byte{
			request.Key: request.Data,
		},
	}
	createdSecret, err := dao.ClientSet.CoreV1().Secrets(request.Namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return createdSecret
}

func (dao *DAO) Update(ctx context.Context, option option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (dao *DAO) Delete(ctx context.Context, option option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (dao *DAO) Close(context context.Context) errors.Error {
	//TODO implement me
	panic("implement me")
}
