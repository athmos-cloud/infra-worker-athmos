package kernel

import "github.com/PaulBarrie/infra-worker/pkg/infrastructure/runtime"

const DefaultWorkdir = "/tmp/infra-worker"
const DefaultRuntime = runtime.DAGGER
const DefaultTmpDir = "/tmp/infra-worker"
const DefaultTerraformImage = "hashicorp/terraform:1.3.9"
