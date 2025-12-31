package controllers

import (
	"net/http"
	"time"

	"sync_golang/config"
	"sync_golang/models"

	"github.com/gin-gonic/gin"
)

func OwnerDashboard(c *gin.Context) {
	today := time.Now().Format("2006-01-02")

	// 1️⃣ TOTAL OMZET SEMUA CABANG
	var totalSales int64
	config.DB.
		Model(&models.Sale{}).
		Where("DATE(created_at) = ?", today).
		Select("COALESCE(SUM(total),0)").
		Scan(&totalSales)

	// 2️⃣ TOTAL TRANSAKSI
	var totalTransactions int64
	config.DB.
		Model(&models.Sale{}).
		Where("DATE(created_at) = ?", today).
		Count(&totalTransactions)

	// 3️⃣ TOTAL KAS MASUK
	var totalCash int64
	config.DB.
		Model(&models.CashTransaction{}).
		Where("type = 'IN' AND DATE(created_at) = ?", today).
		Select("COALESCE(SUM(amount),0)").
		Scan(&totalCash)

	// 4️⃣ OMZET PER CABANG
	var branchSales []struct {
		BranchID uint
		Name     string
		Total    int64
	}
	config.DB.
		Table("sales").
		Select(`
			sales.branch_id,
			branches.name,
			SUM(sales.total) as total
		`).
		Joins("JOIN branches ON branches.id = sales.branch_id").
		Where("DATE(sales.created_at) = ?", today).
		Group("sales.branch_id, branches.name").
		Scan(&branchSales)

	// 5️⃣ PRODUK PALING LAKU
	var topProducts []struct {
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
		Joins("JOIN products ON products.id = sale_items.product_id").
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Where("DATE(sales.created_at) = ?", today).
		Group("sale_items.product_id, products.name").
		Order("total DESC").
		Limit(5).
		Scan(&topProducts)

	c.JSON(http.StatusOK, gin.H{
		"date":               today,
		"total_sales":        totalSales,
		"total_transactions": totalTransactions,
		"total_cash":         totalCash,
		"branch_sales":       branchSales,
		"top_products":       topProducts,
	})
}
