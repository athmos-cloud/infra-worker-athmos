package helm

import (
	"context"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/repo"
	"reflect"
	"sync"
)

var Instance *Repository

var lock = &sync.Mutex{}

type Repository struct {
	HelmClient helmclient.Client
}

func Connect() *Repository {
	lock.Lock()
	defer lock.Unlock()
	if reflect.DeepEqual(Instance, &Repository{}) {
		client, err := helmclient.New(&helmclient.Options{})
		if err != nil {
			logger.Error.Printf("Error creating helm client :  %v", err)
			panic(err)
		}
		err = client.AddOrUpdateChartRepo(
			repo.Entry{
				Name:     "plugins",
				URL:      config.Current.Plugins.Crossplane.Registry.Address,
				Username: config.Current.Plugins.Crossplane.Registry.Username,
				Password: config.Current.Plugins.Crossplane.Registry.Password,
			})
		if err != nil {
			logger.Error.Printf("Error connecting to artifact repository:  %v", err)
			panic(err)
		}
		Instance = &Repository{
			HelmClient: client,
		}
	}
	return Instance
}

// Get retrieves a Helm chart from the Helm repository
// The argument must be a GetHelmChartRequest{chartName, chartVersion}
// The return value is a map[string]interface{} representing the values of the Helm chart
func (r *Repository) Get(_ context.Context, optn option.Option) (interface{}, errors.Error) {
	if optn = optn.SetType(reflect.TypeOf(GetHelmChartRequest{}).String()); !optn.Validate() {
		return nil, errors.InvalidArgument.WithMessage("Argument must be a GetHelmChartRequest{chartName, chartVersion}")
	}
	args := optn.Get().(GetHelmChartRequest)
	logger.Info.Printf("args: %v", args)
	if r.HelmClient == nil {
		logger.Error.Printf("Error loading chart :  %v", args)
	}
	chart, _, err := r.HelmClient.GetChart(
		args.ChartName,
		&action.ChartPathOptions{
			Version:  args.ChartVersion,
			RepoURL:  config.Current.Plugins.Crossplane.Registry.Address,
			Username: config.Current.Plugins.Crossplane.Registry.Username,
			Password: config.Current.Plugins.Crossplane.Registry.Password,
		},
	)
	if err != nil {
		logger.Error.Printf("Error loading chart :  %v", err)
		return nil, errors.ExternalServiceError.WithMessage(
			fmt.Sprintf("Error loading chart :  %v", err),
		)
	}
	return chart, errors.OK
}

func (r *Repository) Create(ctx context.Context, option option.Option) (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) GetAll(ctx context.Context, option option.Option) ([]interface{}, errors.Error) {
	//TODO implement me
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

func init() {
	Instance = &Repository{}
}
