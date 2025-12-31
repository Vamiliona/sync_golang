package controllers

import (
	"net/http"
	"time"

	"sync_golang/config"
	"sync_golang/models"

	"github.com/gin-gonic/gin"
)

func BranchDashboard(c *gin.Context) {
	branchID := c.Param("branch_id")

	today := time.Now().Format("2006-01-02")

	// 1️⃣ TOTAL PENJUALAN HARI INI
	var totalSales int64
	config.DB.
		Model(&models.Sale{}).
		Where("branch_id = ? AND DATE(created_at) = ?", branchID, today).
		Select("COALESCE(SUM(total),0)").
		Scan(&totalSales)

	// 2️⃣ JUMLAH TRANSAKSI
	var totalTransactions int64
	config.DB.
		Model(&models.Sale{}).
		Where("branch_id = ? AND DATE(created_at) = ?", branchID, today).
		Count(&totalTransactions)

	// 3️⃣ TOTAL UANG KAS MASUK
	var cashIn int64
	config.DB.
		Model(&models.CashTransaction{}).
		Where("branch_id = ? AND type = 'IN' AND DATE(created_at) = ?", branchID, today).
		Select("COALESCE(SUM(amount),0)").
		Scan(&cashIn)

	// 4️⃣ STOK PRODUK
	var stocks []struct {
		ProductID uint
		Name      string
		Quantity  int64
	}
	config.DB.
		Table("stocks").
		Select("stocks.product_id, products.name, stocks.quantity").
		Joins("JOIN products ON products.id = stocks.product_id").
		Where("stocks.branch_id = ?", branchID).
		Scan(&stocks)

	// 5️⃣ PENJUALAN PER PRODUK
	var productSales []struct {
		ProductID uint
		Name      string
		Qty       int64
		Total     int64
	}
	config.DB.
		Table("sale_items").
		Select(`
			sale_items.product_id,
			products.name,
			SUM(sale_items.qty) as qty,
			SUM(sale_items.subtotal) as total
		`).
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Joins("JOIN products ON products.id = sale_items.product_id").
		Where("sales.branch_id = ? AND DATE(sales.created_at) = ?", branchID, today).
		Group("sale_items.product_id, products.name").
		Scan(&productSales)

	c.JSON(http.StatusOK, gin.H{
		"date":               today,
		"branch_id":          branchID,
		"total_sales":        totalSales,
		"total_transactions": totalTransactions,
		"cash_in":            cashIn,
		"stocks":             stocks,
		"product_sales":      productSales,
	})
}
