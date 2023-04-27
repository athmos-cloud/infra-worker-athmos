package types

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type DiskMode string

const (
	ReadOnly  DiskMode = "READ_ONLY"
	ReadWrite DiskMode = "READ_WRITE"
)

func DiskModeType(input string) (DiskMode, errors.Error) {
	switch input {
	case "READ_ONLY":
		return ReadOnly, errors.OK
	case "READ_WRITE":
		return ReadWrite, errors.OK
	default:
		return "", errors.InvalidArgument.WithMessage(fmt.Sprintf("invalid disk mode %s", input))
	}
}
