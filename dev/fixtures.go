package dev

import (
	mongo2 "github.com/PaulBarrie/infra-worker/pkg/repository/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/repository/repository/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/context"
)

func apply() {
	mongo2.Client.Create(context.CurrentContext, mongo2.CreateRequestPayload{
		Plugin: mongo.Plugin{
			Name: "test",

		}
	})
}
