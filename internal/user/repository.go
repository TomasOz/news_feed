package user

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(id uint) (*User, error)
	GetByUsername(username string) (*User, error)
	Create(username, passwordHash string) (*User, error)
// 	Update(user *User) error
// 	Delete(id uint) error
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) GetByID(id uint)(*User, error) {
	var user User

	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) GetByUsername(username string) (*User, error) {
	var user User
	
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) Create(username, passwordHash string) (*User, error) {
	user := User{
		Username: username,
		PasswordHash: passwordHash,
	}
	
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}