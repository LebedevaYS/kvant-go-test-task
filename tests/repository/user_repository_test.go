package repository_test

import (
	"errors"
	"go-test-task/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil { // Явная проверка на nil
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	testUser := &models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Age:      25,
	}

	// Правильная настройка мока
	mockRepo.On("CreateUser", testUser).Return(testUser, nil)

	resultUser, err := mockRepo.CreateUser(testUser)

	assert.NoError(t, err)
	assert.Equal(t, testUser.Email, resultUser.Email)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepo)
	testUser := &models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Age:      25,
	}

	// Явно указываем возврат nil и ошибки
	mockRepo.On("CreateUser", testUser).Return(nil, errors.New("email already exists"))

	resultUser, err := mockRepo.CreateUser(testUser)

	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Equal(t, "email already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func (m *MockUserRepo) GetUser(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestGetUser_Success(t *testing.T) {
	// 1. Настраиваем мок
	mockRepo := new(MockUserRepo)
	expectedUser := &models.User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
	}

	// 2. Задаём ожидание:
	// "При вызове GetUser с ID=1 вернуть expectedUser"
	mockRepo.On("GetUser", uint(1)).Return(expectedUser, nil)

	// 3. Вызываем метод
	user, err := mockRepo.GetUser(1)

	// 4. Проверяем
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepo)

	// Ожидаем ошибку при запросе несуществующего пользователя
	mockRepo.On("GetUser", uint(999)).Return(nil, errors.New("not found"))

	user, err := mockRepo.GetUser(999)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "not found", err.Error())
}

func (m *MockUserRepo) UpdateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil { // Явная проверка на nil
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	user := &models.User{ID: 4, Name: "NewName"}

	mockRepo.On("UpdateUser", user).Return(nil, nil)

	user, err := mockRepo.UpdateUser(user)
	assert.NoError(t, err)
	assert.Nil(t, user)
}
