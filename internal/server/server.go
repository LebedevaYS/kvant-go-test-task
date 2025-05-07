package server

import (
	"go-test-task/internal/handlers"
	"go-test-task/internal/middleware"
	"go-test-task/internal/repository"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router       *gin.Engine
	userHandler  *handlers.UserHandler
	orderHandler *handlers.OrderHandler
	authHandler  *handlers.AuthHandler
}

func NewServer(userRepo repository.UserRepository, orderRepo repository.OrderRepository) *Server {
	r := gin.Default()

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userRepo)
	orderHandler := handlers.NewOrderHandler(orderRepo)
	authHandler := handlers.NewAuthHandler(userRepo)

	s := &Server{
		router:       r,
		userHandler:  userHandler,
		orderHandler: orderHandler,
		authHandler:  authHandler,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Public routes
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	s.router.POST("/users", s.userHandler.CreateUser)
	s.router.POST("/auth/login", s.authHandler.Login)

	// Protected routes
	protected := s.router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/users", s.userHandler.GetUsers)
		protected.GET("/users/:id", s.userHandler.GetUser)
		protected.PUT("/users/:id", s.userHandler.UpdateUser)
		protected.DELETE("/users/:id", s.userHandler.DeleteUser)
		protected.POST("/users/:user_id/orders", s.orderHandler.CreateOrder)
		protected.GET("/users/:id/orders", s.orderHandler.GetUserOrders)
	}
}

func (s *Server) Run() error {
	return s.router.Run(":8080")
}
