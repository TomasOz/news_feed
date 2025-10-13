package follow

import (
	"testing"
	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFollowRepository struct {
	mock.Mock
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

func TestFollowService_Follow (t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockFollowRepository)
		wantErr        bool
		wantErrMsg     string
		follower_id	   uint
		followee_id	   uint
		
	}{
		{
			name:     "successful follow other user",
			follower_id: 1,
			followee_id: 2,
			mockSetup: func(m *MockFollowRepository) {
				m.On("AlreadyFollowing", uint(1), uint(2)).Return(false, nil)
				m.On("Follow", uint(1), uint(2)).Return(nil) 
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:     "You can Not follow yourself",
			follower_id: 1,
			followee_id: 1,
			mockSetup: func(m *MockFollowRepository) {},
			wantErr:    true,
			wantErrMsg: "you can not follow yourself",
		},
		{
			name:     "you already follow this user",
			follower_id: 1,
			followee_id: 2,
			mockSetup: func(m *MockFollowRepository) {
				m.On("AlreadyFollowing", uint(1), uint(2)).Return(true, nil)
			},
			wantErr:    true,
			wantErrMsg: "you already follow user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockFollowRepository{}
			service := NewFollowService(mockRepo)
			tt.mockSetup(mockRepo)

			// Act
			err := service.Follow(tt.follower_id, tt.followee_id)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
			
			// Verify all mock expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestFollowService_UnFollow (t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockFollowRepository)
		wantErr        bool
		wantErrMsg     string
		follower_id	   uint
		followee_id	   uint
		
	}{
		{
			name:     "successful unfollow other user",
			follower_id: 1,
			followee_id: 2,
			mockSetup: func(m *MockFollowRepository) {
				m.On("UnFollow", uint(1), uint(2)).Return(nil) 
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:     "You can Not unfollow yourself",
			follower_id: 1,
			followee_id: 1,
			mockSetup: func(m *MockFollowRepository) {},
			wantErr:    true,
			wantErrMsg: "you can not unfollow yourself",
		},
		{
			name:     "User was not found",
			follower_id: 1,
			followee_id: 8,
			mockSetup: func(m *MockFollowRepository) {
				m.On("UnFollow", uint(1), uint(8)).Return(errors.New("user was not found"))
			},
			wantErr:    true,
			wantErrMsg: "user was not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockFollowRepository{}
			service := NewFollowService(mockRepo)
			tt.mockSetup(mockRepo)

			// Act
			err := service.UnFollow(tt.follower_id, tt.followee_id)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
			
			// Verify all mock expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}