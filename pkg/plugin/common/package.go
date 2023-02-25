package common

import (
	"github.com/PaulBarrie/infra-worker/pkg/kernel/types/workdir"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/common/pipeline"
)

type Package struct {
	Workdir  workdir.Workdir
	Pipeline pipeline.Pipeline
}
