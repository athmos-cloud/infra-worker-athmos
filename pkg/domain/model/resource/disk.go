package resource

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
