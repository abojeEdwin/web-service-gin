package service

import (
	"errors"
	"example/web-service-gin/app/models"
	"example/web-service-gin/app/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.repo.CreateUser(user)
}

func (s *UserService) Login(user models.User) (string, error) {
	dbUser, err := s.repo.GetUserByUsername(user.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// In a real application, you would generate a JWT token here.
	// For simplicity, we'll just return a success message or a dummy token.
	return "dummy-token", nil
}
