package user

import (
	"testing"
	// "errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(id uint) (*User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) Create(username, passwordHash string) (*User, error) {
	args := m.Called(username, passwordHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}


func createTestUserWithHashedPassword(username, password string) *User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &User{
		Model: gorm.Model{ID: 1},
		Username: username,
		PasswordHash: string(hashedPassword),
	}
}


func TestUserService_Login(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		password    string
		mockSetup   func(*MockUserRepository)
		wantErr     bool
		wantErrMsg  string
		wantUser    *User
	}{
		{
			name:     "successful login with valid credentials",
			username: "testuser",
			password: "password123",
			mockSetup: func(m *MockUserRepository) {
				user := createTestUserWithHashedPassword("testuser", "password123")
				m.On("GetByUsername", "testuser").Return(user, nil)
			},
			wantErr:    false,
			wantErrMsg: "",
			wantUser:   &User{Username: "testuser"},
		},
		{
			name:     "user not found in database",
			username: "nonexistent",
			password: "password123",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetByUsername", "nonexistent").Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr:    true,
			wantErrMsg: "invalid credentials",
			wantUser:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockUserRepository{}
			service := NewUserService(mockRepo)
			tt.mockSetup(mockRepo)

			// Act
			user, err := service.Login(tt.username, tt.password)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.wantUser.Username, user.Username)
				assert.NotZero(t, user.ID)
			}
			
			// Verify all mock expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}



func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		password    string
		mockSetup   func(*MockUserRepository)
		wantErr     bool
		wantErrMsg  string
		wantUser    *User
	}{
		{
			name:     "successful register with valid credentials",
			username: "newuser",
			password: "password123",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetByUsername", "newuser").Return(nil, gorm.ErrRecordNotFound)
				user := createTestUserWithHashedPassword("newuser", "password123")
				m.On("Create", "newuser", mock.AnythingOfType("string")).Return(user, nil)
			},
			wantErr:    false,
			wantErrMsg: "",
			wantUser:   &User{Username: "newuser"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockUserRepository{}
			service := NewUserService(mockRepo)
			tt.mockSetup(mockRepo)

			// Act
			user, err := service.Register(tt.username, tt.password)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.wantUser.Username, user.Username)
				assert.NotZero(t, user.ID)
			}
			
			// Verify all mock expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}