package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetSecretKey() (string, error) {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return "", errors.New("SECRET_KEY não encontrada no .env")
	}

	return secret, nil
}

func GenerateToken(email, role string, userId int64) (string, error) {
	secretKey, err := GetSecretKey()

	if err != nil {
		return "", errors.New("SECRET_KEY não encontrada no .env")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"role":   role,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}
