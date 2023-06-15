package gcp

import "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"

func toGCPDiskType(diskType instance.DiskType) string {
	switch diskType {
	case instance.DiskTypeSSD:
		return "pd-ssd"
	default:
		return "pd-standard"
	}
}

func fromGCPDiskType(diskType string) instance.DiskType {
	switch diskType {
	case "pd-ssd":
		return instance.DiskTypeSSD
	default:
		return instance.DiskTypeHDD
	}
}

func toGCPDiskMode(diskMode instance.DiskMode) string {
	switch diskMode {
	case instance.DiskModeReadWrite:
		return "READ_WRITE"
	default:
		return "READ_ONLY"
	}
}

func fromGCPDiskMode(diskMode string) instance.DiskMode {
	switch diskMode {
	case "READ_WRITE":
		return instance.DiskModeReadWrite
	default:
		return instance.DiskModeReadOnly
	}
}
