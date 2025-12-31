package models

import "time"

type Product struct {
	ID        uint      `gorm:"primaryKey"`
	StoreID   uint
	Name      string
	Price     int64
	CreatedAt time.Time
}
