package metadata

type Metadata struct {
	Managed bool              `json:"managed"`
	Tags    map[string]string `json:"tags,omitempty"`
}
