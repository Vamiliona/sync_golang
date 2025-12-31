package models

import "time"

type CashTransaction struct {
	ID        uint      `gorm:"primaryKey"`
	StoreID   uint
	BranchID  uint
	SaleID    uint
	Type      string // IN / OUT
	Amount    int64
	CreatedAt time.Time
}
