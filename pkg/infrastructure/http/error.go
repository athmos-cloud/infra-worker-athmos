package http

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/gin-gonic/gin"
	"reflect"
)

func handleError(c *gin.Context, err any) {
	if reflect.TypeOf(err) == reflect.TypeOf(errors.Error{}) {
		errorRaised := err.(errors.Error)
		c.JSON(errorRaised.Code, gin.H{
			"message": errorRaised.ToString(),
		})
	} else {
		c.JSON(500, gin.H{
			"message": err,
		})
	}
}
