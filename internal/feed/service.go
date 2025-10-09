package feed

import (
	"errors"
	"strconv"
	"context"

	"news-feed/internal/post"
	"news-feed/internal/follow"
	"news-feed/internal/cache"
)

type FeedService interface {
	GetFeed(userID uint, limit, offset int) ([]post.Post,  error)
}

type DefaultFeedService struct {
	followRepo follow.FollowRepository
	postRepo   post.PostRepository
	cache 	   *cache.RedisCache
}

func NewFeedService(followRepo follow.FollowRepository, postRepo post.PostRepository, cache *cache.RedisCache) FeedService {
	return &DefaultFeedService{
		followRepo: followRepo, 
		postRepo: postRepo,
		cache: cache,
	}
}

func (s *DefaultFeedService) GetFeed(userID uint, limit, offset int) ([]post.Post, error) {
	ctx := context.Background()
	
	// Lets take post ids from cache first
	postIDs, err := s.getPostIDsFromCache(ctx, userID, limit, offset)

	if err != nil || len(postIDs) == 0 {
		feed, err := s.GetFeedFromDatabase(userID, limit, offset)
		s.addFeedToCache(ctx, userID, feed)

		return feed, err
	}
	
	var postIDsUint []uint
	for _, idStr := range postIDs {
		if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
			postIDsUint = append(postIDsUint, uint(id))
		}
	}
	
	posts, err := s.postRepo.GetPostsByIDs(postIDsUint)

	if err != nil {
		return s.GetFeedFromDatabase(userID, limit, offset)
	}
	
	return posts, nil
}

func (s *DefaultFeedService) getPostIDsFromCache (ctx context.Context, userID uint, limit, offset int) ([]string, error) {
	feedKey := cache.FeedKey(userID)
	
	//By default Offset is 10
	start := int64(offset)
	end := start + int64(limit) - 1
	
	postIDs, err := s.cache.LRange(ctx, feedKey, start, end)
	if err != nil {
		return nil, err
	}
	
	return postIDs, nil
}

func (s *DefaultFeedService) addFeedToCache (ctx context.Context, userID uint, posts []post.Post) (error) {
	feedKey := cache.FeedKey(userID)

	if len(posts) == 0 {
		return nil
	}	

	values := make([]any, len(posts))
	for i, p := range posts {
		values[i] = p.ID
	}

	if err := s.cache.LPush(ctx, feedKey, values); err != nil {
		return err
	}

	return nil
}


func (s *DefaultFeedService) GetFeedFromDatabase(userID uint, limit, offset int) ([]post.Post, error) {
		followeesID, err := s.followRepo.GetFolloweesID(userID)

		if len(followeesID) == 0 {
			return []post.Post{}, nil
		}
	
		if err != nil {
			return nil, errors.New("Internal Error")
		}
		
		feed, err := s.postRepo.GetPostsByUserID(followeesID, limit, offset)
	
		if err != nil {
			return nil, err
		}
	
		return feed, nil
}
