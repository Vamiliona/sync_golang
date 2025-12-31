package routes

import (
	"sync_golang/controllers"
	"sync_golang/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	// ================= PUBLIC =================
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.Refresh)
	r.POST("/logout", controllers.Logout)

	// ================= AUTH BASIC =================
	r.GET("/me", middleware.Auth(), controllers.Me)
	r.POST("/change-password", middleware.Auth(), controllers.ChangePassword)

	// ================= OWNER =================
	r.POST("/stores", middleware.Auth("OWNER"), controllers.CreateStore)
	r.POST("/branches", middleware.Auth("OWNER"), controllers.CreateBranch)
	r.POST("/products", middleware.Auth("OWNER"), controllers.CreateProduct)
	r.POST("/stocks", middleware.Auth("OWNER"), controllers.SetStock)

	// ================= STOCK =================
	r.GET("/stocks", middleware.Auth(), controllers.GetAllStock)
	r.GET("/stocks/:branch_id", middleware.Auth(), controllers.GetStockByBranch)

	// ================= SALES =================
	r.POST("/sales",
		middleware.Auth("KASIR", "OWNER"),
		controllers.CreateSale,
	)

	// ================= DASHBOARD =================
	r.GET("/dashboard/branch/:branch_id",
		middleware.Auth("KASIR", "OWNER"),
		controllers.BranchDashboard,
	)

	r.GET("/dashboard/owner",
		middleware.Auth("OWNER"),
		controllers.OwnerDashboard,
	)

	// ================= CLOSING CASH =================
	r.POST("/cash/close",
		middleware.Auth("KASIR", "OWNER"),
		controllers.CloseCash,
	)
}
