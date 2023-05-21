package mongo

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"sync"
)

var lock = &sync.Mutex{}

var isInit = false

func init() {
	lock.Lock()
	defer lock.Unlock()
	if !isInit {
		logger.Info.Printf("Init mongo client")
		conf := config.Current.Mongo
		uri := "mongodb://" + conf.Username + ":" + conf.Password + "@" + conf.Address + ":" + strconv.Itoa(conf.Port)
		if err := mgm.SetDefaultConfig(nil, conf.Database, options.Client().ApplyURI(uri)); err != nil {
			panic(errors.InternalError.WithMessage(err.Error()))
		}
		isInit = true
	}
}
