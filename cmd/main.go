package main

import (
	"fmt"
	"log"
	"os"

	"go-test-task/internal/config"
	"go-test-task/internal/repository"
	"go-test-task/internal/server"
)

func main() {
	config.Init()

	// Подключение к БД
	db, err := repository.NewDB(buildDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Инициализация сервера
	srv := server.NewServer(userRepo, orderRepo)

	log.Println("Server is running on :8080")
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func buildDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "postgres")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
