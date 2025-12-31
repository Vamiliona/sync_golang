package models

type Permission struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}
