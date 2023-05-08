package helm

import (
	"context"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/ghodss/yaml"
	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
	"os"
	"reflect"
	"sync"
)

const (
	CleanupOnFail    = true
	RepositoryCache  = "/tmp/.helmcache"
	RepositoryConfig = "/tmp/.helmconfig"
	DefaultNamespace = "default"
	pluginsRepoName  = "plugins"
)

var ReleaseClient *ReleaseDAO
var lock = &sync.Mutex{}

type ReleaseDAO struct {
	HelmClient helmclient.Client
	Namespace  string
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	if ReleaseClient == nil {
		logger.Info.Printf("Init helm client...")
		cli := client()
		ReleaseClient = cli
	}
}

func client() *ReleaseDAO {
	kubeConfig, errFile := os.ReadFile(config.Current.Kubernetes.ConfigPath)
	if errFile != nil {
		panic(errors.InternalError.WithMessage(fmt.Sprintf("Error reading kube config file: %v", errFile)))
	}
	cli, err := helmclient.NewClientFromKubeConf(
		&helmclient.KubeConfClientOptions{
			KubeContext: "",
			KubeConfig:  kubeConfig,
			Options: &helmclient.Options{
				RepositoryCache:  RepositoryCache,
				RepositoryConfig: RepositoryConfig,
				Debug:            config.Current.Kubernetes.Helm.Debug,
				Namespace:        DefaultNamespace,
			},
		})
	if err != nil || cli == nil {
		panic(errors.ExternalServiceError.WithMessage(err))
	}
	err = cli.AddOrUpdateChartRepo(
		repo.Entry{
			Name:     pluginsRepoName,
			URL:      config.Current.Plugins.Crossplane.Registry.Address,
			Username: config.Current.Plugins.Crossplane.Registry.Username,
			Password: config.Current.Plugins.Crossplane.Registry.Password,
		})
	if err != nil {
		logger.Info.Printf("Error adding plugins repo: %v", err)
		logger.Info.Println(config.Current.Plugins.Crossplane.Registry.Address, config.Current.Plugins.Crossplane.Registry.Username)
		panic(errors.ExternalServiceError.WithMessage(err.Error()))
	}
	helmClient := &ReleaseDAO{
		HelmClient: cli,
	}
	return helmClient
}

// Get retrieves a Helm chart from the Helm dao
// The argument must be a GetHelmReleaseRequest{chartName, chartVersion}
// The return value is a map[string]interface{} representing the values of the Helm chart
func (r *ReleaseDAO) Get(_ context.Context, optn option.Option) interface{} {
	if optn = optn.SetType(reflect.TypeOf(GetHelmReleaseRequest{}).String()); !optn.Validate() {
		panic(errors.InvalidArgument.WithMessage("Argument must be a GetHelmReleaseRequest{chartName, chartVersion}"))
	}
	args := optn.Get().(GetHelmReleaseRequest)
	release, err := r.HelmClient.GetRelease(args.ReleaseName)
	if err != nil {
		panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error getting release : %v", err)))
	}
	return GetHelmReleaseResponse{release}
}

func (r *ReleaseDAO) Exists(ctx context.Context, o option.Option) (bool, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReleaseDAO) Create(ctx context.Context, request option.Option) interface{} {
	if request = request.SetType(reflect.TypeOf(CreateHelmReleaseRequest{}).String()); !request.Validate() {
		panic(errors.InvalidArgument.WithMessage("Argument must be a GetHelmReleaseRequest{chartName, chartVersion}"))
	}
	args := request.Get().(CreateHelmReleaseRequest)
	yamlBytes, err1 := yaml.Marshal(args.Values)
	if err1 != nil {
		panic(errors.InternalError.WithMessage(err1))
	}

	release, err2 := r.HelmClient.InstallChart(
		ctx,
		&helmclient.ChartSpec{
			ChartName:     fmt.Sprintf("%s/%s", pluginsRepoName, args.ChartName),
			ReleaseName:   args.ReleaseName,
			Version:       args.ChartVersion,
			ValuesYaml:    string(yamlBytes),
			Namespace:     args.Namespace,
			CleanupOnFail: CleanupOnFail,
		},
		&helmclient.GenericHelmOptions{},
	)
	if err2 != nil {
		panic(errors.ExternalServiceError.WithMessage(
			fmt.Sprintf("Error installing chart :  %v", err2),
		))
	}

	return CreateHelmReleaseResponse{Release: release}
}

func (r *ReleaseDAO) Update(_ context.Context, request option.Option) {
	if request = request.SetType(reflect.TypeOf(CreateHelmReleaseRequest{}).String()); !request.Validate() {
		panic(errors.InvalidArgument.WithMessage("Argument must be a UpdateHelmReleaseRequest{ReleaseName, Chart, Namespace, Values}"))
	}
	args := request.Get().(CreateHelmReleaseRequest)
	yamlBytes, err1 := yaml.Marshal(args.Values)
	if err1 != nil {
		panic(errors.InternalError.WithMessage(err1))
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
		panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error upgrading helm release : %v", err)))
	}
}

func (r *ReleaseDAO) Delete(_ context.Context, request option.Option) {
	if request = request.SetType(reflect.TypeOf(DeleteHelmReleaseRequest{}).String()); !request.Validate() {
		panic(errors.InvalidArgument.WithMessage("Argument must be a DeleteHelmReleaseRequest{ReleaseName}"))
	}
	args := request.Get().(DeleteHelmReleaseRequest)
	if err := r.HelmClient.UninstallReleaseByName(args.ReleaseName); err != nil {
		panic(errors.ExternalServiceError.WithMessage(fmt.Sprintf("Error uninstalling helm chart: %v", err)))
	}
}

func (r *ReleaseDAO) GetAll(_ context.Context, _ option.Option) interface{} {
	//TODO implement me
	panic("implement me")
}

func (r *ReleaseDAO) Close(_ context.Context) {
	//TODO implement me
	panic("implement me")
}
