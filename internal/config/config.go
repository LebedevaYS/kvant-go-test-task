package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	SecretKey []byte
)

func Init() {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	// Инициализация секретного ключа для JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET is not set in .env file")
	}
	SecretKey = []byte(secret)
}
