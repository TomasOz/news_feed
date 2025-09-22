package user

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username     string `gorm:"type:varchar(191);uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
}
