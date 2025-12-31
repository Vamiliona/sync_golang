package controllers

import (
	"net/http"

	"sync_golang/config"
	"sync_golang/models"

	"github.com/gin-gonic/gin"
)

func CreateSale(c *gin.Context) {
	var input struct {
		StoreID  uint `json:"store_id"`
		BranchID uint `json:"branch_id"`
		Items []struct {
			ProductID uint `json:"product_id"`
			Qty       int64 `json:"qty"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user := c.MustGet("user").(models.User)

	tx := config.DB.Begin()

	var total int64 = 0
	var saleItems []models.SaleItem

	// 1️⃣ HITUNG TOTAL + VALIDASI STOK
	for _, item := range input.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(404, gin.H{"error": "Product not found"})
			return
		}

		var stock models.Stock
		tx.Where("branch_id = ? AND product_id = ?", input.BranchID, item.ProductID).
			First(&stock)

		if stock.Quantity < item.Qty {
			tx.Rollback()
			c.JSON(400, gin.H{"error": "Stok tidak cukup"})
			return
		}

		subtotal := product.Price * item.Qty
		total += subtotal

		saleItems = append(saleItems, models.SaleItem{
			ProductID: item.ProductID,
			Qty:       item.Qty,
			Price:     product.Price,
			Subtotal:  subtotal,
		})

		// kurangi stok
		tx.Model(&stock).Update("quantity", stock.Quantity-item.Qty)
	}

	// 2️⃣ SIMPAN SALE
	sale := models.Sale{
		StoreID:  input.StoreID,
		BranchID: input.BranchID,
		UserID:   user.ID,
		Total:    total,
	}
	tx.Create(&sale)

	// 3️⃣ SIMPAN SALE ITEMS
	for i := range saleItems {
		saleItems[i].SaleID = sale.ID
		tx.Create(&saleItems[i])
	}

	// 4️⃣ TAMBAH UANG KE KAS
	cash := models.CashTransaction{
		StoreID:  input.StoreID,
		BranchID: input.BranchID,
		SaleID:   sale.ID,
		Type:     "IN",
		Amount:   total,
	}
	tx.Create(&cash)

	tx.Commit()

	c.JSON(200, gin.H{
		"message": "Sale success",
		"sale_id": sale.ID,
		"total":   total,
	})
}
