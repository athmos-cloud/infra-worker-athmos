package http

//func (server *Server) WithResourceController() *Server {
//	server.Engine.GET("/:projectId/:providerID", func(c *gin.Context) {
//		err := errors.OK
//		defer func() {
//			if r := recover(); r != nil {
//				handleError(c, r)
//			}
//		}()
//		resp := server.ResourceService.GetResource(c, resource.GetResourceRequest{
//			ProjectID: c.Param("projectId"),
//			ResourceID: identifier.FromPayload(identifier.IdPayload{
//				Provider: c.Param("providerID"),
//				VPC:      c.Query("vpcID"),
//				Network:  c.Query("networkID"),
//				Subnet:   c.Query("subnetID"),
//				VM:       c.Query("vmID"),
//				Firewall: c.Query("firewallID"),
//			}),
//		})
//
//		c.JSON(err.Code, gin.H{
//			"message": resp.Resource,
//		})
//	})
//
//	return server
//}
