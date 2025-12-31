package models

import "time"

type Stock struct {
	ID        uint
	ProductID uint
	BranchID  uint
	Quantity  int64   // ⬅️ INI INT64
	CreatedAt time.Time
	UpdatedAt time.Time
}
