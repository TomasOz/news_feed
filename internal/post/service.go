package post

import (
	"errors"
)

type PostService interface {
	GetPostByID(id uint) (*Post, error)
	Create(userID uint, title, body string) (*Post, error)
	GetPosts() (*[]Post, error)
}

type DefaultPostService struct {
	repo PostRepository
}

func NewPostService(repo PostRepository) PostService {
	return &DefaultPostService{repo: repo}
}

func (p DefaultPostService) GetPostByID(id uint) (*Post, error) {
	post, err := p.repo.GetPostByID(id)

	if err != nil {
		return nil, errors.New("post not found")
	}

	return post, nil
}

func (p DefaultPostService) GetPosts() (*[]Post, error) {
	post, err := p.repo.GetPosts()

	if err != nil {
		return nil, errors.New("post not found")
	}

	return post, nil
}


func (p DefaultPostService) Create(userID uint, title, body string) (*Post, error) {
	post, err := p.repo.Create(userID, title, body)

	if err != nil {
		return nil, errors.New("post was not created")
	}

	return post, nil
}