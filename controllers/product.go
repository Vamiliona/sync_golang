package controllers

import (
	"sync_golang/config"
	"sync_golang/models"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var input struct {
		StoreID uint
		Name    string
		Price   int64
	}
	c.ShouldBindJSON(&input)

	product := models.Product{
		StoreID: input.StoreID,
		Name:    input.Name,
		Price:   input.Price,
	}

	config.DB.Create(&product)
	c.JSON(200, product)
}
