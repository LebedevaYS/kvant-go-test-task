package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-test-task/internal/models"
	"strconv"
)

func main() {
	// Подключение к БД
	dsn := buildDSN()
	log.Println("Connecting to DB with DSN:", dsn) // Добавьте лог для отладки

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Автомиграции
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	r := gin.Default()

	// Инициализация маршрутов
	setupRoutes(r, db)

	log.Println("Server is running on :8080")
	if err := r.Run(":8080"); err != nil {
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

func setupRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/users", func(c *gin.Context) {
		var input struct {
			Name     string `json:"name" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Age      int    `json:"age" binding:"required,gte=0"`
			Password string `json:"password" binding:"required,min=8"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		user := models.User{
			Name:     input.Name,
			Email:    input.Email,
			Age:      input.Age,
			Password: input.Password, // В реальном приложении нужно хешировать!
		}

		if result := db.Create(&user); result.Error != nil {
			c.JSON(400, gin.H{"error": "User with this email already exists"})
			return
		}

		c.JSON(201, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"age":   user.Age,
		})
	})

	r.GET("/users", func(c *gin.Context) {
		// Параметры запроса с значениями по умолчанию
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		minAge, _ := strconv.Atoi(c.Query("min_age"))
		maxAge, _ := strconv.Atoi(c.Query("max_age"))

		// Рассчитываем offset
		offset := (page - 1) * limit

		// Строим запрос
		query := db.Model(&models.User{})

		// Фильтрация по возрасту
		if minAge > 0 {
			query = query.Where("age >= ?", minAge)
		}
		if maxAge > 0 {
			query = query.Where("age <= ?", maxAge)
		}

		// Получаем общее количество
		var total int64
		query.Count(&total)

		// Получаем данные с пагинацией
		var users []models.User
		query.Offset(offset).Limit(limit).Find(&users)

		// Формируем ответ
		response := gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"users": users,
		}

		c.JSON(200, response)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		// Получаем ID из URL
		id := c.Param("id")

		// Ищем пользователя в базе данных по ID
		var user models.User
		if result := db.First(&user, id); result.Error != nil {
			// Если пользователь не найден, возвращаем 404
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		// Если пользователь найден, возвращаем его данные
		c.JSON(200, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"age":   user.Age,
		})
	})

	// PUT для обновления пользователя
	r.PUT("/users/:id", func(c *gin.Context) {
		// Получаем id пользователя из URL
		id := c.Param("id")

		// Входные данные
		var input struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Age   int    `json:"age"`
		}

		// Привязываем тело запроса к структуре input
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Ищем пользователя по ID
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		// Обновляем данные пользователя
		user.Name = input.Name
		user.Email = input.Email
		user.Age = input.Age

		// Сохраняем изменения в базе данных
		if err := db.Save(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update user"})
			return
		}

		// Возвращаем обновленные данные пользователя
		c.JSON(200, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"age":   user.Age,
		})
	})

	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		if err := db.Unscoped().Delete(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete user"})
			return
		}

		c.Status(204)
	})

}
