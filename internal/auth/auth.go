package auth

import (
	"errors"
	"go-test-task/internal/config"
	"go-test-task/internal/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,                               // ID пользователя
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Срок действия - 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.SecretKey)
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return config.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
