package controllers

import (
	"net/http"
	"time"

	"sync_golang/config"
	"sync_golang/models"

	"github.com/gin-gonic/gin"
)

func CloseCash(c *gin.Context) {
	var input struct {
		StoreID    uint  `json:"store_id"`
		BranchID   uint  `json:"branch_id"`
		OpenCash   int64 `json:"open_cash"`
		ActualCash int64 `json:"actual_cash"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user := c.MustGet("user").(models.User)
	today := time.Now().Format("2006-01-02")

	// 1️⃣ HITUNG TOTAL PENJUALAN HARI INI
	var totalSales int64
	config.DB.
		Model(&models.Sale{}).
		Where("branch_id = ? AND DATE(created_at) = ?", input.BranchID, today).
		Select("COALESCE(SUM(total),0)").
		Scan(&totalSales)

	// 2️⃣ HITUNG KAS SEHARUSNYA
	expectedCash := input.OpenCash + totalSales

	// 3️⃣ SELISIH
	difference := input.ActualCash - expectedCash

	// 4️⃣ SIMPAN CLOSING
	closing := models.CashClosing{
		StoreID:      input.StoreID,
		BranchID:     input.BranchID,
		UserID:       user.ID,
		OpenCash:     input.OpenCash,
		TotalSales:   totalSales,
		ExpectedCash: expectedCash,
		ActualCash:   input.ActualCash,
		Difference:   difference,
		ClosedAt:     time.Now(),
		CreatedAt:    time.Now(),
	}

	config.DB.Create(&closing)

	c.JSON(http.StatusOK, gin.H{
		"message": "Closing kas berhasil",
		"data":    closing,
	})
}
