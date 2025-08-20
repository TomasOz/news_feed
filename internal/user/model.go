package user

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"not null" json:"-"`
}
