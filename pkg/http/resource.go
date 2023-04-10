package http

import (
	"encoding/json"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/common"
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithResourceController() *Server {
	server.Router.GET("/:projectId/:provider/:resourceType/:resourceId", func(c *gin.Context) {
		resp, err := server.ResourceService.GetResource(c, dto.GetResourceRequest{
			ProjectID:    c.Param("projectId"),
			Provider:     common.ProviderType(c.Param("provider")),
			ResourceType: common.ResourceType(c.Param("resourceType")),
			ResourceID:   c.Param("resourceId"),
		})
		if !err.IsOk() {
			c.JSON(err.Code, gin.H{
				"message": err.ToString(),
			})
			return
		}
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
