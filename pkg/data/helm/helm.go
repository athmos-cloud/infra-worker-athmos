package helm

const (
	initialReleaseNumber = "1.0.0"
)

type ReleaseReference struct {
	Name    string `bson:"name"`
	Version string `bson:"version"`
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
