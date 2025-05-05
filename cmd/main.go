package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-test-task/internal/models"

	"crypto/sha256"
	"encoding/hex"
)

func main() {
	// Загрузка переменных окружения
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err)
	}

	// Автоматическая миграция таблиц
	db.AutoMigrate(&models.User{}, &models.Order{})

	// Создание Gin-сервера
	r := gin.Default()

	// Простой эндпоинт для проверки работы сервера
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Сервер успешно запущен и работает!",
		})
	})

	// Роутер для создания пользователя (оставляем ваш старый код)
	r.POST("/users", func(c *gin.Context) {
		// Входная структура (точно соответствующая ТЗ)
		type CreateUserRequest struct {
			Name     string `json:"name" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Age      int    `json:"age" binding:"required,gte=0"`
			Password string `json:"password" binding:"required,min=8"`
		}

		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request data: " + err.Error()})
			return
		}

		// Хеширование пароля (но не возвращаем его в ответе)
		hash := sha256.Sum256([]byte(req.Password))
		passwordHash := hex.EncodeToString(hash[:])

		// Проверка существующего пользователя
		var existingUser models.User
		if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			c.JSON(400, gin.H{"error": "User with this email already exists"})
			return
		}

		// Создание пользователя
		newUser := models.User{
			Name:         req.Name,
			Email:        req.Email,
			Age:          req.Age,
			PasswordHash: passwordHash,
		}

		if err := db.Create(&newUser).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create user"})
			return
		}

		// Ответ строго по ТЗ
		c.JSON(201, gin.H{
			"id":    newUser.ID,
			"name":  newUser.Name,
			"email": newUser.Email,
			"age":   newUser.Age,
		})
	})

	var users = []gin.H{
		{
			"id":    1,
			"name":  "John Doe",
			"email": "john@example.com",
			"age":   30,
		},
		{
			"id":    2,
			"name":  "Jane Smith",
			"email": "jane@example.com",
			"age":   25,
		},
	}

	// Добавьте этот код после инициализации маршрутизатора (r := gin.Default())
	r.GET("/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"users": users,
			"count": len(users),
		})
	})

	// Запуск сервера с логом
	log.Printf("🚀 Сервер запущен и доступен на http://localhost:%s", port)
	log.Printf("🔍 Проверить статус можно по /health")
	r.Run(":" + port)
}
