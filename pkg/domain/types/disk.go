package types

type DiskMode string

const (
	ReadOnly  DiskMode = "READ_ONLY"
	ReadWrite DiskMode = "READ_WRITE"
)

func DiskModeType(input string) DiskMode {
	switch input {
	case "READ_ONLY":
		return ReadOnly
	case "READ_WRITE":
		return ReadWrite
	default:
		return ""
	}
}
