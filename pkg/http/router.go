package http

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.New()
	router.GET("/projects/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/projects", func(c *gin.Context) {
		c.JSON(201, gin.H{
			"id": "1",
		})
	})
	router.GET("/projects/owner/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/plugins/:providerType/:resourceType", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/resources/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := router.Run(fmt.Sprintf(":%d", config.Current.Http.Port))
	if err != nil {
		panic(err)
	}
}
