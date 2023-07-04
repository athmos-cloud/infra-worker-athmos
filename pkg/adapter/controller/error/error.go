package errorCtrl

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/gin-gonic/gin"
	"reflect"
)

func RaiseError(ctx context.Context, err any) {
	logger.Error.Println(err)
	if reflect.TypeOf(err) == reflect.TypeOf(errors.Error{}) {
		errs := err.(errors.Error)
		ctx.JSON(errs.Code, gin.H{"message": errs.Message})
	} else {
		ctx.JSON(500, gin.H{"message": errors.InternalError.WithMessage(fmt.Sprintf("%v", err)).Message})
	}
}
