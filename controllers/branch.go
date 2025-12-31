package controllers

import (
	"sync_golang/config"
	"sync_golang/models"

	"github.com/gin-gonic/gin"
)

func CreateBranch(c *gin.Context) {
	var input struct {
		StoreID uint
		Name    string
		Address string
	}
	c.ShouldBindJSON(&input)

	branch := models.Branch{
		StoreID: input.StoreID,
		Name:    input.Name,
		Address: input.Address,
	}

	config.DB.Create(&branch)
	c.JSON(200, branch)
}
