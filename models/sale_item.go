package models

type SaleItem struct {
	ID        uint `gorm:"primaryKey"`
	SaleID    uint
	ProductID uint
	Qty       int64
	Price     int64
	Subtotal  int64
}
