package models

import (
	"gorm.io/gorm"
)

// User представляет сущность пользователя в базе данных
type User struct {
	gorm.Model        // включает ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `json:"name"`
	Email      string `json:"email" gorm:"uniqueIndex"`
	Age        int    `json:"age"`
	Password   string `json:"-"` // не возвращается в JSON-ответах
}
