package instance

type VMDisk struct {
	Type       DiskType `json:"type"`
	Mode       DiskMode `json:"mode"`
	SizeGib    int      `json:"size_gib"`
	AutoDelete bool     `json:"auto_delete"`
}

type VMDiskList []VMDisk

type DiskMode string

const (
	DiskModeReadWrite DiskMode = "READ_WRITE"
	DiskModeReadOnly  DiskMode = "READ_ONLY"
)

type DiskType string

const (
	DiskTypeSSD DiskType = "SSD"
	DiskTypeHDD DiskType = "HDD"
)

type SqlDbDisk struct {
	Type               DiskType `json:"type"`
	SizeGib            int      `json:"size_gib"`
	Autoresize         bool     `json:"autoresize,omitempty" default:"true"`
	AutoresizeLimitGib int      `json:"autoresize_limit_gib,omitempty" default:"0"`
}
