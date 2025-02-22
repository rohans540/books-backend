package models

import "gorm.io/gorm"

type Book struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `gorm:"not null" json:"title"`
	Author string `gorm:"not null" json:"author"`
	Year   int    `json:"year"`
}

func MigrateBooks(db *gorm.DB) {
	db.AutoMigrate(&Book{})
}
