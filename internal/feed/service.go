package feed

import (
	"errors"
	"strconv"
	"strings"
	"fmt"
	"time"

	"news-feed/internal/post"
	"news-feed/internal/follow"
)

type FeedService interface {
	GetFeed(userID uint, limit int, cursor string) ([]post.Post, string, error)
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

func (s *DefaultFeedService) GetFeed(userID uint, limit int, cursor string) ([]post.Post, string, error) {
	followeesID, err := s.followRepo.GetFolloweesID(userID)

	if len(followeesID) == 0 {
		return []post.Post{},  "", nil
	}

	if err != nil {
		return nil, "", errors.New("Internal Error")
	}

    var createdAt time.Time
    var lastID uint
    if cursor != "" {
		parts := strings.Split(cursor, ",")
		createdAt, _ = time.Parse(time.RFC3339, parts[0])
		u64, _ := strconv.ParseUint(parts[1], 10, 64)
		lastID = uint(u64)
    }

	feed, err := s.postRepo.GetPostsByUserID(followeesID, limit, createdAt, lastID)

	if err != nil {
        return nil, "", err
    }

    var nextCursor string
    if len(feed) == limit {
        last := feed[len(feed)-1]
        nextCursor = fmt.Sprintf("%s,%d", last.CreatedAt.Format(time.RFC3339), last.ID)
    }

    return feed, nextCursor, nil
}