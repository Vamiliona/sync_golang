package models

import "time"

type CashSession struct {
	ID         uint      `gorm:"primaryKey"`
	StoreID    uint
	BranchID   uint
	UserID     uint
	OpenCash   int64
	ActualCash int64
	Expected   int64
	Difference int64
	ClosedAt   time.Time
	CreatedAt  time.Time
}
