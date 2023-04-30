package http

import (
	"encoding/json"
	"fmt"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithResourceController() *Server {
	server.Router.GET("/:projectId/:provider/:resourceType/:resourceId", func(c *gin.Context) {
		err := errors.OK
		defer func() {
			if r := recover(); r != nil {
				err = r.(errors.Error)
				c.JSON(err.Code, gin.H{
					"message": err.ToString(),
				})
			}
		}()
		resp := server.ResourceService.GetResource(c, dto.GetResourceRequest{
			ProjectID:  c.Param("projectId"),
			ResourceID: c.Param("resourceId"),
		})

		jsonBytes, errMarshal := json.Marshal(resp)
		if errMarshal != nil {
			c.JSON(500, gin.H{
				"message": fmt.Sprintf("Error marshalling response: %s", errMarshal),
			})
			return
		}
		c.JSON(err.Code, gin.H{
			"message": string(jsonBytes[:]),
		})
	})

	return server
}
