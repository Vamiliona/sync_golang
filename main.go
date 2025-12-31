package main

import (
	"os"

	"sync_golang/config"
	"sync_golang/models"
	"sync_golang/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Mode release kalau di server
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Connect Database
	config.ConnectDB()

	// Auto migrate tables
	config.DB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Store{},
		&models.Branch{},
		&models.Product{},
		&models.Stock{},
		&models.Sale{},
		&models.SaleItem{},
		&models.CashSession{},
	)

	// Register routes
	routes.RegisterRoutes(r)

	// Ambil PORT dari ENV (Railway / Fly / Docker)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // local default
	}

	// Run server
	r.Run(":" + port)
}
