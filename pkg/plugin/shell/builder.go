package shell

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/types/workdir"
	"github.com/PaulBarrie/infra-worker/pkg/plugin"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/common"
)

type Builder struct {
}

func (b *Builder) Load(ctx context.Context, option option.Option) (common.Plugin, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (b *Builder) Vendor(ctx context.Context, plugin common.Plugin) (workdir.Workdir, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (b *Builder) Build(ctx context.Context, payload plugin.BuildPluginPayload) (workdir.Workdir, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (b *Builder) Register(ctx context.Context, location common.Location) (common.Plugin, errors.Error) {
	//TODO implement me
	panic("implement me")
}

