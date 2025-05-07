package repository

import (
	"go-test-task/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetUserOrders(userID uint) ([]models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetUserOrders(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}
