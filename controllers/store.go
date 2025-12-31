package controllers

import (
	"sync_golang/config"
	"sync_golang/models"

	"github.com/gin-gonic/gin"
)

func CreateStore(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var input struct {
		Name string
		Type string
	}
	c.ShouldBindJSON(&input)

	store := models.Store{
		Name:    input.Name,
		Type:    input.Type,
		OwnerID: user.ID,
	}

	config.DB.Create(&store)
	c.JSON(200, store)
}
