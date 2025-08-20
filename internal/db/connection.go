package db

import (
	"news-feed/internal/user"

	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB


func InitDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatal("Failed to auto-migrate:", err)
	}

	return db
}