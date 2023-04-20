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

func (r *DAO) Get(ctx context.Context, option option.Option) (interface{}, errors.Error) {
	if option = option.SetType(reflect.TypeOf(GetResourceRequest{}).String()); !option.Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			"Argument must be a kubernetes.GetResourceRequest{resourceId, namespace, name}",
		)
	}
	request := option.Get().(GetResourceRequest)

	resource, err := r.DynamicClient.Resource(request.ResourceID).
		Namespace(request.Namespace).Get(ctx, request.Name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.NotFound.WithMessage(err.Error())
	}

	return resource, errors.OK
}

func (r *DAO) Exists(ctx context.Context, o option.Option) (bool, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (r *DAO) GetAll(ctx context.Context, option option.Option) (interface{}, errors.Error) {
	if option = option.SetType(reflect.TypeOf(GetListResourceRequest{}).String()); !option.Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			"Argument must be a kubernetes.GetListResourceRequest{resourceId, namespace, labels}",
		)
	}
	request := option.Get().(GetListResourceRequest)
	var options metav1.ListOptions
	if request.Labels != nil {
		options = metav1.ListOptions{
			LabelSelector: labelsToString(request.Labels),
		}
	}
	list, err := r.DynamicClient.Resource(request.ResourceID).Namespace(request.Namespace).List(ctx, options)
	if err != nil {
		return nil, errors.NotFound.WithMessage(err.Error())
	}
	return list.Items, errors.OK
}

// Create namspace from a given name string
func (r *DAO) Create(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	if optn = optn.SetType(reflect.TypeOf(CreateNamespaceRequest{}).String()); !optn.Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			"Argument must be a kubernetes.CreateNamespaceRequest{name}",
		)
	}
	request := optn.Get().(CreateNamespaceRequest)
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: request.Name,
		},
	}
	namespace, err := r.ClientSet.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.ExternalServiceError.WithMessage(err.Error())
	}
	return namespace, errors.OK
}

func (r *DAO) Update(ctx context.Context, option option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (r *DAO) Delete(ctx context.Context, option option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (r *DAO) Close(context context.Context) errors.Error {
	//TODO implement me
	panic("implement me")
}
