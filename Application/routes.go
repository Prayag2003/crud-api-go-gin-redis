package application

import (
	handler "github.com/Prayag2003/crud-api-golang/Handler"
	order "github.com/Prayag2003/crud-api-golang/Repository/Order"
	"github.com/gin-gonic/gin"
)

func (a *App) loadRoutes() {
	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, Gin!"})
	})

	orderGroup := router.Group("/orders")
	{
		orderHandler := &handler.Order{
			Repo: &order.RedisRepo{
				Client: a.rdb,
			},
		}
		orderGroup.POST("/", orderHandler.Create)
		orderGroup.GET("/", orderHandler.List)
		orderGroup.GET("/:id", orderHandler.OrderById)
		orderGroup.PUT("/:id", orderHandler.UpdateById)
		orderGroup.DELETE("/:id", orderHandler.DeleteById)
	}
	a.router = router
}
