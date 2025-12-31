package controllers

import (
	"net/http"
	"sync_golang/config"
	"sync_golang/models"
	"time"

	"github.com/gin-gonic/gin"
)

func SetStock(c *gin.Context) {
	var input struct {
		ProductID uint `json:"product_id"`
		BranchID  uint `json:"branch_id"`
		Quantity  int64  `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var stock models.Stock
	// cek apakah sudah ada stok untuk product di branch
	err := config.DB.Where("product_id = ? AND branch_id = ?", input.ProductID, input.BranchID).
		First(&stock).Error

	if err != nil {
		// belum ada -> create
		newStock := models.Stock{
			ProductID: input.ProductID,
			BranchID:  input.BranchID,
			Quantity:  input.Quantity,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		config.DB.Create(&newStock)
		c.JSON(http.StatusOK, newStock)
		return
	}

	// sudah ada -> update
	stock.Quantity = input.Quantity
	stock.UpdatedAt = time.Now()
	config.DB.Save(&stock)
	c.JSON(http.StatusOK, stock)
}

// GetAllStock -> ambil semua stok
func GetAllStock(c *gin.Context) {
	var stocks []models.Stock
	config.DB.Find(&stocks)
	c.JSON(http.StatusOK, stocks)
}

// GetStockByBranch -> ambil stok per cabang
func GetStockByBranch(c *gin.Context) {
	branchID := c.Param("branch_id")

	var stocks []models.Stock
	if err := config.DB.Where("branch_id = ?", branchID).Find(&stocks).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	c.JSON(http.StatusOK, stocks)
}
