package application

import (
	handler "github.com/Prayag2003/go-microservice/Handler"
	"github.com/gin-gonic/gin"
)

func loadRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, Gin!"})
	})

	orderGroup := router.Group("/orders")
	{
		orderHandler := &handler.Order{}
		orderGroup.POST("/", orderHandler.Create)
		orderGroup.GET("/", orderHandler.List)
		orderGroup.GET("/:id", orderHandler.OrderById)
		orderGroup.PUT("/:id", orderHandler.UpdateById)
		orderGroup.DELETE("/:id", orderHandler.DeleteById)
	}
	return router
}
