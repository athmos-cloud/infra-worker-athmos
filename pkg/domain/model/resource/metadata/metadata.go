package metadata

type Metadata struct {
	Status  Status            `json:"status"`
	Managed bool              `json:"managed"`
	Tags    map[string]string `json:"tags,omitempty"`
}
