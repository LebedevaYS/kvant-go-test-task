package handlers

import (
	"net/http"
	"strconv"

	"go-test-task/internal/models"
	"go-test-task/internal/repository"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	repo repository.OrderRepository
}

func NewOrderHandler(repo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input struct {
		Product  string  `json:"product" binding:"required"`
		Quantity int     `json:"quantity" binding:"required,min=1"`
		Price    float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := models.Order{
		UserID:   uint(userID),
		Product:  input.Product,
		Quantity: input.Quantity,
		Price:    input.Price,
	}

	if err := h.repo.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	orders, err := h.repo.GetUserOrders(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
