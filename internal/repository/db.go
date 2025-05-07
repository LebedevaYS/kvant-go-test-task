package repository

import (
	"go-test-task/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Автомиграции
	if err := db.AutoMigrate(&models.User{}, &models.Order{}); err != nil {
		return nil, err
	}

	// Добавление внешнего ключа
	db.Exec(`
		ALTER TABLE orders 
		ADD CONSTRAINT fk_orders_user 
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
	`)

	return db, nil
}
