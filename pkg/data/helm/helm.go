package helm

import "github.com/kamva/mgm/v3"

const (
	initialReleaseNumber = "1.0.0"
)

type ReleaseReference struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name"`
	Version          string `bson:"version"`
}

func NewRelease(name string) ReleaseReference {
	return ReleaseReference{
		Name:    name,
		Version: initialReleaseNumber,
	}
}

func (release *ReleaseReference) Equals(other ReleaseReference) bool {
	return release.Name == other.Name && release.Version == other.Version
}
