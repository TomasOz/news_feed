package post

import (
	"testing"
	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockPostRepository struct {
	mock.Mock
}

func createTestPost(id, user_id uint, title, body string) *Post {
	return &Post{
		Model: gorm.Model{ID: id},
		UserID: user_id,
		Title: title,
		Body: body,
	}
}

func (m *MockPostRepository) GetPostByID(id uint) (*Post, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*Post), args.Error(1)
}

func (m *MockPostRepository) Create(userID uint, title, body string) (*Post, error) {
    args := m.Called(userID, title, body)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*Post), args.Error(1)
}

func (m *MockPostRepository) GetPosts() (*[]Post, error) {
    args := m.Called()
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*[]Post), args.Error(1)
}

func (m *MockPostRepository) GetPostsByUserID(followeesID []uint, limit, offset int) ([]Post, error) {
    args := m.Called(followeesID, limit, offset)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]Post), args.Error(1)
}

func (m *MockPostRepository) GetPostsByIDs(ids []uint) ([]Post, error) {
    args := m.Called(ids)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]Post), args.Error(1)
}

func TestPostService_GetPostByID_Success(t *testing.T) {
	postMockRepository := &MockPostRepository{}
	service := NewPostService(postMockRepository)
	expectedPost := createTestPost(1, 2, "TestPost", "TestBody")
	postMockRepository.On("GetPostByID", uint(1)).Return(expectedPost, nil)

	post, err := service.GetPostByID(1)

	assert.NoError(t, err)
    assert.Equal(t, uint(1), post.ID)
    postMockRepository.AssertExpectations(t)
}

func TestPostService_GetPostByID_NotFound(t *testing.T) {
	postMockRepository := &MockPostRepository{}
	service := NewPostService(postMockRepository)
	postMockRepository.On("GetPostByID", uint(1)).Return(nil, gorm.ErrRecordNotFound)

	post, err := service.GetPostByID(1)

	assert.Error(t, err)
    assert.Nil(t, post)
    postMockRepository.AssertExpectations(t)
}

func TestPostService_GetPostByID_DatabaseError(t *testing.T) {
    postMockRepository := &MockPostRepository{}
    service := NewPostService(postMockRepository)
    postMockRepository.On("GetPostByID", uint(1)).Return(nil, errors.New("database error"))

    result, err := service.GetPostByID(1)

    assert.Error(t, err)
    assert.Nil(t, result)
    postMockRepository.AssertExpectations(t)
}