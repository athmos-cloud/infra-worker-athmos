package helm

import "helm.sh/helm/v3/pkg/release"

type GetHelmReleaseRequest struct {
	ReleaseName string
}

type GetHelmReleaseResponse struct {
	Release *release.Release
}

type CreateHelmReleaseRequest struct {
	ChartName    string
	ChartVersion string
	ReleaseName  string
	Values       map[string]interface{}
	Namespace    string
}

type UpdateHelmReleaseRequest struct {
	ReleaseName string
	ChartName   string
	Namespace   string
	Values      map[string]interface{}
}

type DeleteHelmReleaseRequest struct {
	ReleaseName string
}
