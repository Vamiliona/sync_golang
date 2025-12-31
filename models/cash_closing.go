package models

import "time"

type CashClosing struct {
	ID           uint
	StoreID      uint
	BranchID     uint
	UserID       uint
	OpenCash     int64
	TotalSales   int64
	ExpectedCash int64
	ActualCash   int64
	Difference   int64
	ClosedAt     time.Time
	CreatedAt    time.Time
}
