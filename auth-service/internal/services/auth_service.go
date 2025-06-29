package services


import (
	"auth-service/internal/models"
	"auth-service/internal/pkg/jwt"
	"auth-service/internal/repository"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    Register(login string, email string, password string) (*models.User, *models.Tokens, error)
    Authorize(email string, password string) (*models.User, *models.Tokens, error, bool)
    Refresh(refreshToken string) (string, string, error)
    ValidateToken(tokenValue string) (bool, string)
    Logout(refreshToken string) (bool, error)
}

type AuthServiceImpl struct {
	userRepo *repository.UserRepository
	tokenRepo *repository.TokenRepository
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository) AuthService {
	return &AuthServiceImpl{userRepo: userRepo, tokenRepo: tokenRepo}
}

func (s *AuthServiceImpl) Register(login string, email string, password string) (*models.User, *models.Tokens, error) {
	hashedPswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}
	user:=&models.User{
		Email: email,
		Password: string(hashedPswd),
		Login: login,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	accesstoken, refreshtoken, err := jwt.GenerateJWTRefreshTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}
	refreshToken := &models.RefreshToken{
		UserID: user.ID,
		Token:  refreshtoken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.tokenRepo.Create(refreshToken); err != nil {
		return nil, nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	tokens := &models.Tokens{
		AccessToken: accesstoken,
		RefreshToken: refreshtoken,
	}

	return user, tokens, nil
}

func (s *AuthServiceImpl) Authorize(email string, password string) (*models.User, *models.Tokens, error, bool) {
	
	user, err, exists := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, nil, err, exists
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		//реализовать решение лучше
		return nil, nil, fmt.Errorf("Invalid password"), exists
	}

	accesstoken, refreshtoken, err := jwt.GenerateJWTRefreshTokens(user.ID)
	if err != nil {
		return nil, nil, err, exists
	}
	refreshToken := &models.RefreshToken{
		UserID: user.ID,
		Token:  refreshtoken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.tokenRepo.Create(refreshToken); err != nil {
		return nil, nil, fmt.Errorf("failed to save refresh token: %w", err), exists
	}

	tokens := &models.Tokens{
		AccessToken: accesstoken,
		RefreshToken: refreshtoken,
	}

	return user, tokens, err, exists
}

func (s *AuthServiceImpl) Refresh(refreshToken string) (accesstoken string, refreshtoken string, err error) {
	token, err :=s.tokenRepo.GetToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	if time.Now().After(token.ExpiresAt) {
		return "", "", fmt.Errorf("Token time expired")
	}

	t1, t2, err := jwt.GenerateJWTRefreshTokens(token.UserID)
	if err != nil {
		return "","",err
	}
	return t1, t2, nil
}

func (s *AuthServiceImpl) ValidateToken(tokenValue string) (bool, string) {
	if tokenValue == ""  {
		return false, "Access denied"
	}

	//todo: проверка что это юзерский токен
	// token, err := s.tokenRepo.GetToken(resfreshToken)
	// if err!=nil {
	// 	return false, ""
	// }

	parts:=strings.Split(tokenValue, ".")
	payload := parts[1]

	// payload := tokenValue
	dst := make([]byte, base64.RawURLEncoding.DecodedLen(len(payload)))
	n, err := base64.RawURLEncoding.Decode(dst, []byte(payload))
	payloadJson := dst[:n]
	if err != nil {
		log.Println("decode error:", err)
		return false, "Decode error"
	}

	type Token struct {
			Sub uint `json:"sub"`
			Exp int64 `json:"exp"`
	}
	var token Token
	err = json.Unmarshal(payloadJson, &token)
	if err != nil {
			log.Println("Can't unmarshal")
			return false, "Can't unmarshal"
	}
	log.Println(time.Unix(token.Exp, 0))
	if token.Exp < time.Now().Unix() {
			log.Print("Time ran out")
			return false, "Time ran out"
	}

	return true, ""
}

func (s *AuthServiceImpl) Logout(refreshToken string) (bool, error) {
	success, err := s.tokenRepo.DeleteToken(refreshToken)
	if err != nil {
		return false, err
	}
	return success, nil
}