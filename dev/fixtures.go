package dev

import (
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/context"
)

func apply() {
	mongo.Client.Create(context.CurrentContext, mongo.CreateRequestPayload{
		Plugin: mongo.Plugin{
			Name: "test",

		}
	})
}
