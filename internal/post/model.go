package post

import (
	"gorm.io/gorm"
	"news-feed/internal/user"
)

type Post struct {
	gorm.Model
	
	UserID uint     	`gorm:"not null" json:"user_id"`
	User   user.User 	`gorm:"foreignKey:UserID"`
	Title  string   	`gorm:"not null" json:"title"`
	Body   string		`gorm:"not null" json:"body"`
	
}

