package main

import (
	"fmt"
	"log"
	"os"

	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-test-task/internal/models"

	"go-test-task/internal/config"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"go-test-task/internal/auth"

	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	config.Init()
	// Подключение к БД
	dsn := buildDSN()
	log.Println("Connecting to DB with DSN:", dsn) // Добавьте лог для отладки

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Автомиграции
	if err := db.AutoMigrate(&models.User{}, &models.Order{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	db.Exec(`
		ALTER TABLE orders 
		ADD CONSTRAINT fk_orders_user 
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
	`)

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
	// Защищенные маршруты (для них нужен JWT)
	protected := r.Group("/")
	protected.Use(JWTAuthMiddleware())
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

		// Логируем поступление запроса
		log.Println("Received request to create user:", input.Email)

		if err := c.ShouldBindJSON(&input); err != nil {
			log.Printf("Error binding input: %v\n", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Хеширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v\n", err)
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}

		// Создаем пользователя
		user := models.User{
			Name:     input.Name,
			Email:    input.Email,
			Age:      input.Age,
			Password: string(hashedPassword),
		}

		// Логируем создание пользователя
		log.Printf("Creating user with email: %s\n", input.Email)

		if result := db.Create(&user); result.Error != nil {
			log.Printf("Error creating user: %v\n", result.Error)
			c.JSON(400, gin.H{"error": "User with this email already exists"})
			return
		}

		// Логируем успешное создание пользователя
		log.Printf("User created successfully: %s\n", user.Email)

		c.JSON(201, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"age":   user.Age,
		})
	})

	protected.GET("/users", func(c *gin.Context) {
		// Параметры запроса с значениями по умолчанию
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		minAge, _ := strconv.Atoi(c.Query("min_age"))
		maxAge, _ := strconv.Atoi(c.Query("max_age"))

		log.Printf("Запрос списка пользователей — page: %d, limit: %d, min_age: %d, max_age: %d", page, limit, minAge, maxAge)

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
		result := query.Offset(offset).Limit(limit).Find(&users)

		if result.Error != nil {
			log.Printf("Ошибка при получении пользователей: %v", result.Error)
			c.JSON(500, gin.H{"error": "Ошибка при получении пользователей"})
			return
		}

		log.Printf("Найдено пользователей всего: %d, возвращено: %d", total, len(users))

		// Формируем ответ
		response := gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"users": users,
		}

		c.JSON(200, response)
	})

	protected.GET("/users/:id", func(c *gin.Context) {
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
	protected.PUT("/users/:id", func(c *gin.Context) {
		// Получаем id пользователя из URL
		id := c.Param("id")
		log.Printf("Попытка обновления пользователя с ID: %s", id)

		// Входные данные
		var input struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Age   int    `json:"age"`
		}

		// Привязываем тело запроса к структуре input
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Printf("Ошибка валидации данных при обновлении пользователя ID %s: %v", id, err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Ищем пользователя по ID
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			log.Printf("Пользователь с ID %s не найден", id)
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		// Обновляем данные пользователя
		user.Name = input.Name
		user.Email = input.Email
		user.Age = input.Age

		// Сохраняем изменения в базе данных
		if err := db.Save(&user).Error; err != nil {
			log.Printf("Не удалось обновить пользователя с ID %s: %v", id, err)
			c.JSON(500, gin.H{"error": "Failed to update user"})
			return
		}

		log.Printf("Пользователь с ID %d успешно обновлен: имя=%s, email=%s, возраст=%d", user.ID, user.Name, user.Email, user.Age)

		// Возвращаем обновленные данные пользователя
		c.JSON(200, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"age":   user.Age,
		})
	})

	protected.DELETE("/users/:id", func(c *gin.Context) {
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

	protected.POST("/users/:user_id/orders", func(c *gin.Context) {
		var input struct {
			Product  string  `json:"product" binding:"required"`
			Quantity int     `json:"quantity" binding:"required,min=1"`
			Price    float64 `json:"price" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := db.First(&user, c.Param("user_id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		order := models.Order{
			UserID:   user.ID,
			Product:  input.Product,
			Quantity: input.Quantity,
			Price:    input.Price,
		}

		if err := db.Create(&order).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create order"})
			return
		}

		c.JSON(201, order)
	})

	protected.GET("/users/:id/orders", func(c *gin.Context) {
		var user models.User
		id := c.Param("id")

		// Проверяем, существует ли пользователь
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		// Получаем заказы пользователя
		var orders []models.Order
		if err := db.Where("user_id = ?", id).Find(&orders).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve orders"})
			return
		}

		c.JSON(200, orders)
	})

	r.POST("/auth/login", func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		// Чтение JSON из запроса
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Получаем пользователя из базы данных по email
		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Сравнение пароля с хешированным в базе
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Генерация токена
		token, err := auth.GenerateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		// Возвращаем токен
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Разделяем "Bearer" и сам токен
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be 'Bearer {token}'"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Верифицируем токен
		token, err := auth.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			c.Abort()
			return
		}

		// Устанавливаем claims в контекст
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", claims["sub"])
		}

		c.Next()
	}
}
