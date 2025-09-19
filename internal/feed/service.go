package feed

import (
	"errors"

	"news-feed/internal/post"
	"news-feed/internal/follow"
)

type FeedService interface {
	GetFeed(userID uint) ([]post.Post, error)
}

type DefaultFeedService struct {
	followRepo follow.FollowRepository
	postRepo   post.PostRepository
}

func NewFeedService(followRepo follow.FollowRepository, postRepo post.PostRepository) FeedService {
	return &DefaultFeedService{
		followRepo: followRepo, 
		postRepo: postRepo,
	}
}

func (s *DefaultFeedService) GetFeed(userID uint) ([]post.Post, error) {
	followeesID, err := s.followRepo.GetFolloweesID(userID)

	if len(followeesID) == 0 {
		return []post.Post{}, nil
	}

	if err != nil {
		return nil, errors.New("Internal Error")
	}

	feed, err := s.postRepo.GetPostsByUserID(followeesID)

	return feed, nil
}