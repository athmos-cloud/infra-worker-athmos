package test

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/mongo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func Init(t *testing.T) *gnomock.Container {
	mongoC := InitMongo(t)

	return mongoC
}

func InitMongo(t *testing.T) *gnomock.Container {
	user := "test"
	password := "test"
	db := "infra-test"
	p := mongo.Preset(
		mongo.WithData("./testdata/"),
		mongo.WithUser(user, password),
	)
	c, err := gnomock.Start(p)
	require.NoError(t, err)
	addr := c.DefaultAddress()
	uri := fmt.Sprintf("mongodb://%s:%s@%s", user, password, addr)

	errMgm := mgm.SetDefaultConfig(nil, db, options.Client().ApplyURI(uri))
	require.NoError(t, errMgm)

	return c
}
