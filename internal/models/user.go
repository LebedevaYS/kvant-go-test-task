package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string  `gorm:"size:255;not null" json:"name" binding:"required"`
	Email        string  `gorm:"size:255;uniqueIndex;not null" json:"email" binding:"required,email"`
	Age          int     `gorm:"not null" json:"age" binding:"required,min=1"`
	PasswordHash string  `gorm:"size:255;not null" json:"-"` // Исключаем из JSON
	Orders       []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

// BeforeCreate хук для GORM
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Age < 0 {
		return errors.New("age cannot be negative")
	}
	return nil
}

// SetPassword хеширует пароль и сохраняет хеш
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword проверяет соответствие пароля хешу
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// PublicUser структура для ответа API (без чувствительных данных)
type PublicUser struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `gorm:"not null"` // Храним только хеш
	Age          int    `json:"age"`
}

// ToPublic преобразует User в PublicUser
func (u *User) ToPublic() PublicUser {
	return PublicUser{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Age:          u.Age,
	}
}
