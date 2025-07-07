package jwt

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type JWTFunctional struct {
	secretKey []byte
}

func NewJWTFunctional() JWTFunctional {
	secretKey := GetJWTSecretKey()
	return JWTFunctional{secretKey: secretKey}
}

func GetJWTSecretKey() []byte {
	if err := godotenv.Load(".env.local"); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
	key:=os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		log.Println("wrong key format")
	}
	return []byte(key)
}


func (j *JWTFunctional) GenerateJWTRefreshTokens(userID uint) (accessToken string, refreshToken string, err error){
	var jwtSecretKey = j.secretKey

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
		"exp": jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
	}

	RefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t2, err := RefreshToken.SignedString(jwtSecretKey)
	if err != nil {
        return "", "", err
    }

	return t1, t2, nil
}

func (j *JWTFunctional) GenerateJWTAccessToken(userID uint) (refreshToken string, err error){
	var jwtSecretKey = j.secretKey

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

func  (j *JWTFunctional) VerifyKey(token string) (bool, error){
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        return j.secretKey, nil
    })

	if err != nil {
		return false, err
	}

	return true, nil
}