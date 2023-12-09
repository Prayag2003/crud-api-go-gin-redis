package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Order struct{}

func (o *Order) Create(c *gin.Context) {
	c.String(http.StatusOK, "Inside create method")
}

func (o *Order) List(c *gin.Context) {
	c.String(http.StatusOK, "List method was called")
}

func (o *Order) OrderById(c *gin.Context) {
	c.String(http.StatusOK, "Gets an order by ID")
}

func (o *Order) UpdateById(c *gin.Context) {
	c.String(http.StatusOK, "Updates an order by ID")
}

func (o *Order) DeleteById(c *gin.Context) {
	c.String(http.StatusOK, "Deletes an order by ID")
}
