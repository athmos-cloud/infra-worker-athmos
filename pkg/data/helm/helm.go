package helm

import "github.com/kamva/mgm/v3"

const (
	initialReleaseNumber = "1.0.0"
)

type ReleaseReference struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string   `bson:"name"`
	LatestVersion    string   `bson:"LatestVersion"`
	Versions         []string `bson:"versions,omitempty"` //not handled yet
}

func NewRelease(name string) ReleaseReference {
	return ReleaseReference{
		Name:          name,
		LatestVersion: initialReleaseNumber,
	}
}

func (release *ReleaseReference) Equals(other ReleaseReference) bool {
	return release.Name == other.Name && release.LatestVersion == other.LatestVersion
}
