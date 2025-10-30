package feed

import (
    "testing"
    "context"
    // "errors"
    "news-feed/internal/post"
    "news-feed/internal/cache"

    "github.com/alicebob/miniredis/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "gorm.io/gorm"
    
)

type MockFollowRepository struct {
	mock.Mock
}

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) GetPostByID(id uint) (*post.Post, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*post.Post), args.Error(1)
}

func (m *MockPostRepository) Create(userID uint, title, body string) (*post.Post, error) {
    args := m.Called(userID, title, body)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*post.Post), args.Error(1)
}

func (m *MockPostRepository) GetPosts() (*[]post.Post, error) {
    args := m.Called()
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*[]post.Post), args.Error(1)
}

func (m *MockPostRepository) GetPostsByUserID(followeesID []uint, limit, offset int) ([]post.Post, error) {
    args := m.Called(followeesID, limit, offset)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]post.Post), args.Error(1)
}

func (m *MockPostRepository) GetPostsByIDs(ids []uint) ([]post.Post, error) {
    args := m.Called(ids)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]post.Post), args.Error(1)
}

func (m *MockFollowRepository) Follow(follower_id, followee_id uint) error {
	args := m.Called(follower_id, followee_id)
	return args.Error(0)
}

func (m *MockFollowRepository) UnFollow(follower_id, followee_id uint) error {
	args := m.Called(follower_id, followee_id)
	return args.Error(0) 
}

func (m *MockFollowRepository) AlreadyFollowing(follower_id, followee_id uint) (bool, error) {
	args := m.Called(follower_id, followee_id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockFollowRepository) GetFolloweesID(follower_id uint) ([]uint, error) {
	args := m.Called(follower_id)
	if args.Get(0) == nil {
		return  nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
}

func (m *MockFollowRepository) GetFollowersID(followee_id uint) ([]uint, error) {
	args := m.Called(followee_id)
	if args.Get(0) == nil {
		return  nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
}
//

func TestFeedService_GetFeed_CacheHit(t *testing.T) {
    // Start in-memory Redis
    mr, err := miniredis.Run()
    if err != nil {
        t.Fatalf("failed to start miniredis: %v", err)
    }
    defer mr.Close()

    redisCache := cache.NewRedisCache(mr.Addr(), "", 0)

    postMockRepo := &MockPostRepository{}
    followMockRepo := &MockFollowRepository{}

    service := NewFeedService(followMockRepo, postMockRepo, redisCache)

    ctx := context.Background()

    cacheKey := cache.FeedKey(1)
    _ = redisCache.LPush(ctx, cacheKey, []any{1, 2, 3}...)


    postMockRepo.On("GetPostsByIDs", []uint{3, 2, 1}).Return([]post.Post{
        {Model: gorm.Model{ID: 1}, UserID: 10, Title: "Post 1", Body: "Body 1"},
        {Model: gorm.Model{ID: 2}, UserID: 10, Title: "Post 2", Body: "Body 2"},
        {Model: gorm.Model{ID: 3}, UserID: 10, Title: "Post 3", Body: "Body 3"},
    }, nil)

    // Execute
    posts, err := service.GetFeed(1, 10, 0)

    // Assert
    assert.NoError(t, err)
    assert.Len(t, posts, 3)
    postMockRepo.AssertExpectations(t)
}