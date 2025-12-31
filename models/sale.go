package models

import "time"

type Sale struct {
	ID        uint      `gorm:"primaryKey"`
	StoreID   uint
	BranchID  uint
	UserID    uint
	Total     int64
	CreatedAt time.Time
}
