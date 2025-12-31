package models

import "time"

type Branch struct {
	ID        uint      `gorm:"primaryKey"`
	StoreID   uint
	Name      string
	Address   string
	CreatedAt time.Time
}
