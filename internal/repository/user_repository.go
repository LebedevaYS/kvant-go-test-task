package repository

import (
	"errors"
	"go-test-task/internal/models"

	"gorm.io/gorm"
)

var (
	ErrEmailExists = errors.New("email already exists")
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUser(id uint) (*models.User, error)
	GetUsers(page, limit, minAge, maxAge int) ([]models.User, int64, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUser(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) DeleteUser(id uint) error {
	var user models.User

	// Сначала проверяем существование пользователя
	if err := r.db.First(&user, id).Error; err != nil {
		return err // gorm.ErrRecordNotFound если не найден
	}

	// Мягкое удаление (soft delete) с Unscoped()
	if err := r.db.Unscoped().Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	// Проверяем, существует ли пользователь
	var existingUser models.User
	if err := r.db.First(&existingUser, user.ID).Error; err != nil {
		return err
	}

	// Проверяем, не используется ли email другим пользователем
	if user.Email != existingUser.Email {
		var count int64
		r.db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
		if count > 0 {
			return ErrEmailExists
		}
	}

	return r.db.Model(&models.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
	}).Error
}

func (r *userRepository) GetUsers(page, limit, minAge, maxAge int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	if minAge > 0 {
		query = query.Where("age >= ?", minAge)
	}
	if maxAge > 0 {
		query = query.Where("age <= ?", maxAge)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
