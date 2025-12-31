package models

import "time"

type Store struct {
	ID        uint      `gorm:"primaryKey"`
	OwnerID   uint
	Name      string
	Type      string
	CreatedAt time.Time
}
