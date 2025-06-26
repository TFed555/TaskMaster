package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTRefreshTokens(userID uint) (accessToken string, refreshToken string, err error){
	var jwtSecretKey = []byte("secret-key")

	payload := jwt.MapClaims {
		"sub": userID,
		"exp": jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	}

	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t1, err := AccessToken.SignedString(jwtSecretKey)
	if err != nil {
        return "", "", err
    }

	payload = jwt.MapClaims {
		"sub": userID,
		"exp": jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
	}

	RefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t2, err := RefreshToken.SignedString(jwtSecretKey)
	if err != nil {
        return "", "", err
    }

	return t1, t2, nil
}

func GenerateJWTAccessToken(userID uint) (refreshToken string, err error){
	var jwtSecretKey = []byte("secret-key")

	payload := jwt.MapClaims {
		"sub": userID,
		"exp": jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	}

	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t1, err := AccessToken.SignedString(jwtSecretKey)
	if err != nil {
        return "", err
    }

	return t1, nil
}