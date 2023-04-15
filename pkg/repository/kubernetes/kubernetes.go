package kubernetes

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	RandomKubeConfigUUIDLength = 6
)

type Repository struct {
	Client dynamic.Interface
}

func Connect(ctx context.Context) (*Repository, errors.Error) {
	// Build the configuration object
	r := &Repository{}
	config := ctrl.GetConfigOrDie()
	dynamicCli := dynamic.NewForConfigOrDie(config)
	r.Client = dynamicCli
	return r, errors.OK
}

func (r *Repository) Get(ctx context.Context, option option.Option) (interface{}, errors.Error) {
	if option = option.SetType(reflect.TypeOf(GetResourceRequest{}).String()); !option.Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			"Argument must be a kubernetes.GetResourceRequest{resourceId, namespace, name}",
		)
	}
	request := option.Get().(GetResourceRequest)

	resource, err := r.Client.Resource(request.ResourceID).
		Namespace(request.Namespace).Get(ctx, request.Name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.NotFound.WithMessage(err.Error())
	}

	return resource, errors.OK
}

func (r *Repository) Exists(ctx context.Context, o option.Option) (bool, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) GetAll(ctx context.Context, option option.Option) (interface{}, errors.Error) {
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
	list, err := r.Client.Resource(request.ResourceID).Namespace(request.Namespace).List(ctx, options)
	if err != nil {
		return nil, errors.NotFound.WithMessage(err.Error())
	}
	return list.Items, errors.OK
}

func (r *Repository) Create(ctx context.Context, option option.Option) (interface{}, errors.Error) {
	//TODO implement
	panic("implement me")
}

func (r *Repository) Update(ctx context.Context, option option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Delete(ctx context.Context, option option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Close(context context.Context) errors.Error {
	//TODO implement me
	panic("implement me")
}
