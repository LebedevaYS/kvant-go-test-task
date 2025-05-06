package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),     // из .env
		os.Getenv("DB_USER"),     // postgres
		os.Getenv("DB_PASSWORD"), // postgres
		os.Getenv("DB_NAME"),     // users_orders
		os.Getenv("DB_PORT"),     // 5432
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL!")
	return db, nil
}
