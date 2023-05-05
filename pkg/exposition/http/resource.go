package http

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithResourceController() *Server {
	server.Router.GET("/:projectId/:providerID", func(c *gin.Context) {
		err := errors.OK
		defer func() {
			if r := recover(); r != nil {
				handleError(c, r)
			}
		}()
		resp := server.ResourceService.GetResource(c, resource.GetResourceRequest{
			ProjectID: c.Param("projectId"),
			ResourceID: identifier.Build(identifier.IdPayload{
				ProviderID: c.Param("providerID"),
				VPCID:      c.Query("vpcID"),
				NetworkID:  c.Query("networkID"),
				SubnetID:   c.Query("subnetID"),
				VMID:       c.Query("vmID"),
				FirewallID: c.Query("firewallID"),
			}),
		})

		c.JSON(err.Code, gin.H{
			"message": resp.Resource,
		})
	})

	return server
}
