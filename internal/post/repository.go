package post

import (
	"gorm.io/gorm"

	"time"
)

type PostRepository interface {
	GetPostByID(id uint) (*Post, error)
	Create(userID uint, title, body string) (*Post, error)
	GetPosts() (*[]Post, error)
	GetPostsByUserID(followeesID []uint, limit int, cursorCreatedAt time.Time, cursorID uint) ([]Post, error)

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

func (r *GormPostRepository) GetPostsByUserID(
    followeesID []uint, 
    limit int, 
    cursorCreatedAt time.Time, 
    cursorID uint,
) ([]Post, error) {
	var posts []Post

	query := r.db.
		Where("user_id IN ?", followeesID)
	
	if !cursorCreatedAt.IsZero() && cursorID != 0 {
		query = query.
			Where("(created_at < ?) OR (created_at = ? AND id < ?)",
				cursorCreatedAt, cursorCreatedAt, cursorID)
	}
	
	err := query.
		Order("created_at DESC, id DESC").
		Limit(limit).
		Preload("User").
		Find(&posts).Error
	
	if err != nil {
		return nil, err
	}
	return posts, nil
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


