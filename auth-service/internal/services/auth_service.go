package services

import (
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(email string, password string) (*models.User, error) {
	hashedPswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user:=&models.User{
		Email: email,
		Password: string(hashedPswd),
		Login: "",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Authorize(email string, password string) (bool, error) {
	if user, err := s.userRepo.GetByEmail(email); err != nil {
		return false, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, err
	}
	return true, err
}