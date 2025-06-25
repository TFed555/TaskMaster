package services

import (
	"auth-service/internal/models"
	"auth-service/internal/pkg/jwt"
	"auth-service/internal/repository"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	tokenRepo *repository.TokenRepository
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository) *AuthService {
	return &AuthService{userRepo: userRepo, tokenRepo: tokenRepo}
}

func (s *AuthService) Register(login string, email string, password string) (*models.User, error) {
	hashedPswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user:=&models.User{
		Email: email,
		Password: string(hashedPswd),
		Login: login,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Authorize(email string, password string) (*models.Tokens, error, bool) {
	
	user, err, exists := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err, exists
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		//реализовать решение лучше
		return nil, fmt.Errorf("Invalid password"), exists
	}

	accesstoken, refreshtoken, err := jwt.GenerateJWTRefreshTokens(user.ID)
	if err != nil {
		return nil, err, exists
	}
	refreshToken := &models.RefreshToken{
		UserID: user.ID,
		Token:  refreshtoken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.tokenRepo.Create(refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err), exists
	}

	// err = s.tokenRepo.Create(refreshToken)
	// if err != nil {
	// 	return nil, err, exists
	// }

	tokens := &models.Tokens{
		AccessToken: accesstoken,
		RefreshToken: refreshtoken,
	}

	return tokens, err, exists
}