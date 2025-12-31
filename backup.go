package main

import (
	"sync_golang/config"
	"sync_golang/models"
	"sync_golang/routes"

	"github.com/gin-gonic/gin"
)

func main1() {
	r := gin.Default()

	config.ConnectDB()

	config.DB.AutoMigrate(
	&models.User{},
	&models.RefreshToken{},
	&models.Store{},
	&models.Branch{},
	&models.Product{},
	&models.Stock{},
)


	routes.RegisterRoutes(r)

	r.Run(":8080")
}
