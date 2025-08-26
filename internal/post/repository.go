package post

import (
	"gorm.io/gorm"
)

type PostRepository interface {
	GetPostByID(id uint) (*Post, error)
	Create(userID uint, title, body string) (*Post, error)
	GetPosts() (*[]Post, error)

// 	Update(user *User) error
// 	Delete(id uint) error
}

type GormPostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &GormPostRepository{db: db}
}

func (r *GormPostRepository) GetPostByID(id uint) (*Post, error) {
	var post Post

	if err := r.db.First(&post, id).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *GormPostRepository) GetPosts() (*[]Post, error) {
	var posts []Post

	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}

	return &posts, nil
}

func (r *GormPostRepository) Create(userID uint, title, body string) (*Post, error) {
	post := Post{
		UserID: userID,
		Title: title,
		Body: body,
	}

	if err := r.db.Create(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}


