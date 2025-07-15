package services

import (
	"auth-service/internal/grpc/grpc_client"
	"auth-service/internal/models"
	"auth-service/internal/pkg/domain_models"
	"auth-service/internal/pkg/jwt"
	"auth-service/internal/repository"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type AuthService interface {
    Register(params domain_models.RegisterParams) (int, domain_models.Tokens, error)
    Authorize(params domain_models.AuthorizeParams) (models.User, domain_models.Tokens, error)
    Refresh(refreshToken string) (string, string, error)
    ValidateToken(tokenValue string) (bool, string)
	UpdateAccessToken(refreshToken string) (bool, string, error)
	ParseUserId(token string) (uint, string)
	UpdateUser(domain_models.UpdateUserParams) (int, error)
	DeleteUser(userID uint) (bool, error)
}

type AuthServiceImpl struct {
	userRepo repository.UserRepository
	jwtFunc jwt.JWTFunctional
	mediaClient grpc_client.GRPCMediaService
}

func NewAuthService(userRepo repository.UserRepository, 
			jwtFunc jwt.JWTFunctional, mediaClient grpc_client.GRPCMediaService) AuthService {
	return &AuthServiceImpl{userRepo: userRepo,
			jwtFunc: jwtFunc, mediaClient: mediaClient,}
}

func (s AuthServiceImpl) Register(params domain_models.RegisterParams) (int, domain_models.Tokens, error) {
	hashedPswd, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return -1, domain_models.Tokens{}, err
	}
	user:=&models.User{
		Email: params.Email,
		Password: string(hashedPswd),
		Login: params.Name,
	}

	if err := s.userRepo.Create(user); err != nil {
		return -1, domain_models.Tokens{}, err
	}

	accesstoken, refreshtoken, err := s.jwtFunc.GenerateJWTRefreshTokens(user.ID)
	if err != nil {
		return -1, domain_models.Tokens{}, err
	}

	tokens := domain_models.Tokens{
		AccessToken: accesstoken,
		RefreshToken: refreshtoken,
	}

	return int(user.ID), tokens, nil
}

func (s AuthServiceImpl) Authorize(params domain_models.AuthorizeParams) (models.User, domain_models.Tokens, error) {
	user, err := s.userRepo.GetByEmail(params.Email)
	imgPath := ""
	if user.ImgPath != nil {
		if s.mediaClient == nil {
			mediaCon, err := grpc.Dial("localhost:50050", grpc.WithInsecure())
			if err != nil {
				log.Fatalf("Failed to connect to grpc server: %v", err)
			}
			s.mediaClient = grpc_client.NewGRPCMediaService(mediaCon)
			log.Print("Connected to media-service")
			imgPath = *user.ImgPath
			log.Print("ImgPath ", imgPath)
			defer mediaCon.Close()
			picLink, err := s.mediaClient.GetAvatarPic(*user.ImgPath)
			if err != nil {
				log.Fatalf("Something went wrong; %v", err)
			}
			user.ImgPath = &picLink
		}
	}

	if err != nil {
		return models.User{}, domain_models.Tokens{}, fmt.Errorf("User not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return models.User{}, domain_models.Tokens{}, fmt.Errorf("Invalid password")
	}

	accesstoken, refreshtoken, err := s.jwtFunc.GenerateJWTRefreshTokens(user.ID)
	if err != nil {
		return models.User{}, domain_models.Tokens{}, err
	}

	tokens := domain_models.Tokens{
		AccessToken: accesstoken,
		RefreshToken: refreshtoken,
	}

	log.Printf("Called from auth_service, access_token: %s", tokens.AccessToken)

	return user, tokens, err
}

func (s AuthServiceImpl) Refresh(refreshToken string) (string, string, error) {
	_, errMsg := s.ValidateToken(refreshToken)
	if errMsg != "" {
		return "", "", errors.New("Invalid token")
	}

	userID, errMsg := s.ParseUserId(refreshToken)
	if errMsg != "" {
		return "", "", errors.New("Can't get user ID")
	}

	t1, t2, err := s.jwtFunc.GenerateJWTRefreshTokens(userID)
	if err != nil {
		return "","",err
	}

	return t1, t2, nil
}

func (s AuthServiceImpl) ValidateToken(tokenValue string) (bool, string) {
	if tokenValue == ""  {
		return false, "Access denied"
	}

	verified, err := s.jwtFunc.VerifyKey(tokenValue)
	if !verified || err != nil {
		return false, "Access denied"
	}

	parts:=strings.Split(tokenValue, ".")
	payload := parts[1]

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

	if _, err := s.userRepo.GetUserByID(token.Sub); err != nil {
		log.Print("Not found %d userID", token.Sub)
		return false, "User not found"
	}

	return true, ""
}


func (s AuthServiceImpl) UpdateAccessToken(refreshToken string) (bool, string, error) {
	userID, err:= s.ParseUserId(refreshToken)
	if err != "" {
		return false, "", errors.New("Can't get UserID")
	}
	if refreshToken != "" {
		newValue, err := s.jwtFunc.GenerateJWTAccessToken(userID)
		if err != nil {
			log.Print("Generation failed")
			return false, "", err
		}

		log.Printf("Called from UpdateAccessToken: %s", newValue)
		return true, newValue, nil
	}

	return false, "", errors.New("Can't update access token")
}

func (s AuthServiceImpl) ParseUserId(token string) (uint, string) {
	tokenMas := strings.Split(token,".")
	log.Print(tokenMas)
	payload := tokenMas[1]
	dst := make([]byte, base64.RawURLEncoding.DecodedLen(len(payload)))
	n, err := base64.RawURLEncoding.Decode(dst, []byte(payload))
	payloadJson := dst[:n]
	if err != nil {
		log.Println("decode error:", err)
		return 0, err.Error()
	}

	type Token struct {
			Sub uint `json:"sub"`
			Exp int64 `json:"exp"`
	}
	var newToken Token
	err = json.Unmarshal(payloadJson, &newToken)
	if err != nil {
			log.Println("Can't unmarshal")
			return 0, err.Error()
	}
	return newToken.Sub, ""
}

func (s AuthServiceImpl) UpdateUser(params domain_models.UpdateUserParams) (int, error) {
	userData := domain_models.User{ID: params.ID}
	if params.Email != nil {
        userData.Email = *params.Email
    }
    if params.Name != nil {
        userData.Name = *params.Name
    }
    if params.Password != nil {
        userData.Password = *params.Password
    }
	id, err := s.userRepo.UpdateUser(userData)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s AuthServiceImpl) DeleteUser(userID uint) (bool, error) {
	answer, err := s.userRepo.DeleteUser(userID)
	if err != nil {
		return answer, err
	}
	return answer, nil
}