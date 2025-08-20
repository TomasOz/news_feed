package user

import (
	"golang.org/x/crypto/bcrypt"
	"errors"
)

type UserService interface {
	Login(username, password string) (*User, error)
	Register(username, password string) (*User, error)
}

type DefaultUserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &DefaultUserService{repo: repo}
}

func (u DefaultUserService) Login(username, password string) (*User, error) {
	user, err := u.repo.GetByUsername(username)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (u DefaultUserService) Register(username, password string) (*User, error) {
	_, err := u.repo.GetByUsername(username)

	if err == nil {
		return nil, errors.New("user already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user, err := u.repo.Create(username, string(passwordHash))

	if err != nil {
		return nil, err
	}

	return user, nil
}