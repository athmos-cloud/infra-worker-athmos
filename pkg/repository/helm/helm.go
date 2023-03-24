package helm

import (
	"context"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/ghodss/yaml"
	"github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
	"os"
	"reflect"
)

const (
	CleanupOnFail    = true
	RepositoryCache  = "/tmp/.helmcache"
	RepositoryConfig = "/tmp/.helmconfig"
)

type ReleaseRepository struct {
	HelmClient helmclient.Client
	Namespace  string
}

func Client(namespace string) (*ReleaseRepository, errors.Error) {
	kubeConfig, errFile := os.ReadFile(config.Current.Kubernetes.ConfigPath)
	if errFile != nil {
		msg := fmt.Sprintf("Error reading kube config file: %s", errFile)
		logger.Error.Print(msg)
		return nil, errors.InternalError.WithMessage(msg)
	}
	client, err := helmclient.NewClientFromKubeConf(
		&helmclient.KubeConfClientOptions{
			KubeContext: "",
			KubeConfig:  kubeConfig,
			Options: &helmclient.Options{
				RepositoryCache:  RepositoryCache,
				RepositoryConfig: RepositoryConfig,
				Debug:            config.Current.Kubernetes.Helm.Debug,
				Namespace:        namespace,
			},
		})
	if err != nil || client == nil {
		logger.Error.Printf("Error creating helm client :  %v", err)
		return nil, errors.ExternalServiceError.WithMessage(err)
	}
	err = client.AddOrUpdateChartRepo(
		repo.Entry{
			Name:     "plugins",
			URL:      config.Current.Plugins.Crossplane.Registry.Address,
			Username: config.Current.Plugins.Crossplane.Registry.Username,
			Password: config.Current.Plugins.Crossplane.Registry.Password,
		})
	logger.Info.Print("Plugins repo added")
	if err != nil {
		logger.Error.Printf("Error connecting to artifact repository:  %v", err)
		return nil, errors.ExternalServiceError.WithMessage(err)
	}
	helmClient := &ReleaseRepository{
		HelmClient: client,
	}
	return helmClient, errors.OK
}

// Get retrieves a Helm chart from the Helm repository
// The argument must be a GetHelmReleaseRequest{chartName, chartVersion}
// The return value is a map[string]interface{} representing the values of the Helm chart
func (r *ReleaseRepository) Get(_ context.Context, optn option.Option) (interface{}, errors.Error) {
	if optn = optn.SetType(reflect.TypeOf(GetHelmReleaseRequest{}).String()); !optn.Validate() {
		return nil, errors.InvalidArgument.WithMessage("Argument must be a GetHelmReleaseRequest{chartName, chartVersion}")
	}
	args := optn.Get().(GetHelmReleaseRequest)
	release, err := r.HelmClient.GetRelease(args.ReleaseName)
	if err != nil {
		return nil, errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error getting release : %v", err))
	}
	return GetHelmReleaseResponse{release}, errors.OK
}

func (r *ReleaseRepository) Create(ctx context.Context, request option.Option) (interface{}, errors.Error) {
	if request = request.SetType(reflect.TypeOf(CreateHelmReleaseRequest{}).String()); !request.Validate() {
		return nil, errors.InvalidArgument.WithMessage("Argument must be a GetHelmReleaseRequest{chartName, chartVersion}")
	}
	args := request.Get().(CreateHelmReleaseRequest)
	yamlBytes, err1 := yaml.Marshal(args.Values)
	if err1 != nil {
		return "", errors.InternalError.WithMessage(err1)
	}
	release, err2 := r.HelmClient.InstallChart(
		ctx,
		&helmclient.ChartSpec{
			ChartName:     args.ChartName,
			ReleaseName:   args.ReleaseName,
			Version:       args.ChartVersion,
			ValuesYaml:    string(yamlBytes),
			Namespace:     args.Namespace,
			CleanupOnFail: CleanupOnFail,
		},
		&helmclient.GenericHelmOptions{},
	)
	if err2 != nil {
		logger.Error.Printf("Error installing chart :  %v", err2)
		return nil, errors.ExternalServiceError.WithMessage(
			fmt.Sprintf("Error installing chart :  %v", err2),
		)
	}

	return release, errors.Created
}

func (r *ReleaseRepository) Update(ctx context.Context, request option.Option) errors.Error {
	if request = request.SetType(reflect.TypeOf(CreateHelmReleaseRequest{}).String()); !request.Validate() {
		return errors.InvalidArgument.WithMessage("Argument must be a UpdateHelmReleaseRequest{ReleaseName, ChartName, Namespace, Values}")
	}
	args := request.Get().(CreateHelmReleaseRequest)
	yamlBytes, err1 := yaml.Marshal(args.Values)
	if err1 != nil {
		return errors.InternalError.WithMessage(err1)
	}
	chartSpec := helmclient.ChartSpec{
		ReleaseName: args.ReleaseName,
		ChartName:   args.ChartName,
		Namespace:   "default",
		UpgradeCRDs: true,
		Wait:        true,
		ValuesYaml:  string(yamlBytes),
	}
	if _, err := r.HelmClient.InstallOrUpgradeChart(context.Background(), &chartSpec, nil); err != nil {
		return errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error upgrading helm release : %v", err))
	}
	return errors.Accepted
}

func (r *ReleaseRepository) Delete(ctx context.Context, request option.Option) errors.Error {
	if request = request.SetType(reflect.TypeOf(DeleteHelmReleaseRequest{}).String()); !request.Validate() {
		return errors.InvalidArgument.WithMessage("Argument must be a DeleteHelmReleaseRequest{ReleaseName}")
	}
	args := request.Get().(DeleteHelmReleaseRequest)
	if err := r.HelmClient.UninstallReleaseByName(args.ReleaseName); err != nil {
		return errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error uninstalling helm chart: %v", err))
	}
	return errors.OK
}

func (r *ReleaseRepository) GetAll(ctx_ context.Context, _ option.Option) ([]interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReleaseRepository) Close(context context.Context) errors.Error {
	//TODO implement me
	panic("implement me")
}
