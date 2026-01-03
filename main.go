package main

import (
	"log"
	"os"

	"sync_golang/config"
	"sync_golang/models"
	"sync_golang/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	// ===============================
	// SET GIN MODE
	// ===============================
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ===============================
	// DEBUG ENV DB (PENTING DI RAILWAY)
	// ===============================
	log.Println("MYSQL_HOST:", os.Getenv("MYSQL_HOST"))
	log.Println("MYSQL_PORT:", os.Getenv("MYSQL_PORT"))
	log.Println("MYSQL_USER:", os.Getenv("MYSQL_USER"))
	log.Println("MYSQL_DATABASE:", os.Getenv("MYSQL_DATABASE"))

	// ===============================
	// INIT GIN
	// ===============================
	r := gin.Default()

	// ===============================
	// CONNECT DATABASE
	// ===============================
	config.ConnectDB()

	// ===============================
	// AUTO MIGRATE (INI YANG BUAT TABLE)
	// ===============================
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Store{},
		&models.Branch{},
		&models.Product{},
		&models.Stock{},
		&models.Sale{},
		&models.SaleItem{},
		&models.CashSession{},
		&models.CashTransaction{},
		&models.CashClosing{},
	)

	if err != nil {
		log.Fatal("❌ AutoMigrate failed:", err)
	}

	log.Println("✅ AutoMigrate success")

	// ===============================
	// REGISTER ROUTES
	// ===============================
	routes.RegisterRoutes(r)

	// ===============================
	// PORT (WAJIB UNTUK RAILWAY)
	// ===============================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Server running on port", port)

	// ===============================
	// RUN SERVER
	// ===============================
	r.Run(":" + port)
}
