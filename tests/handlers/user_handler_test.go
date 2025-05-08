package handlers_test

import (
	"errors"
	"go-test-task/internal/handlers"
	"go-test-task/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

// Полная реализация всех методов репозитория
func (m *MockUserRepo) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepo) GetUser(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetUsers(page, limit, minAge, maxAge int) ([]models.User, int64, error) {
	args := m.Called(page, limit, minAge, maxAge)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepo) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Тест успешного создания пользователя
func TestUserHandler_CreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	expectedUser := &models.User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   25,
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).
		Return(expectedUser, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(
		"POST",
		"/users",
		strings.NewReader(`{"name":"Alice","email":"alice@example.com","password":"password123","age":25}`),
	)
	c.Request.Header.Set("Content-Type", "application/json")

	handler := handlers.NewUserHandler(mockRepo)
	handler.CreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"id":1,"name":"Alice","email":"alice@example.com","age":25}`, w.Body.String())
	mockRepo.AssertExpectations(t)
}

// Тест ошибки при создании пользователя
func TestUserHandler_CreateUser_Error(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			"Duplicate email",
			`{"name":"Alice","email":"exists@example.com","password":"pwd","age":25}`,
			errors.New("email exists"),
			http.StatusBadRequest,
			`{"error":"email exists"}`,
		},
		{
			"Invalid input",
			`{"name":"","email":"invalid","password":"pwd","age":0}`,
			nil,
			http.StatusBadRequest,
			`{"error":"Key: 'User.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			if tt.mockError != nil {
				mockRepo.On("CreateUser", mock.Anything).Return(nil, tt.mockError)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/users", strings.NewReader(tt.input))
			c.Request.Header.Set("Content-Type", "application/json")

			handler := handlers.NewUserHandler(mockRepo)
			handler.CreateUser(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockRepo.AssertExpectations(t)
		})
	}
}

// Тест получения пользователя
func TestUserHandler_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	expectedUser := &models.User{
		ID:    1,
		Name:  "Bob",
		Email: "bob@example.com",
		Age:   30,
	}

	mockRepo.On("GetUser", uint(1)).Return(expectedUser, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/users/1", nil)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := handlers.NewUserHandler(mockRepo)
	handler.GetUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"id":1,"name":"Bob","email":"bob@example.com","age":30}`, w.Body.String())
	mockRepo.AssertExpectations(t)
}

// Тест обновления пользователя
func TestUserHandler_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userID := uint(1)
	updatedUser := &models.User{
		ID:    userID,
		Name:  "Bob Updated",
		Email: "bob.updated@example.com",
		Age:   31,
	}

	mockRepo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(nil)
	mockRepo.On("GetUser", userID).Return(updatedUser, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(
		"PUT",
		"/users/1",
		strings.NewReader(`{"name":"Bob Updated","email":"bob.updated@example.com","age":31}`),
	)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request.Header.Set("Content-Type", "application/json")

	handler := handlers.NewUserHandler(mockRepo)
	handler.UpdateUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"id":1,"name":"Bob Updated","email":"bob.updated@example.com","age":31}`, w.Body.String())
	mockRepo.AssertExpectations(t)
}

// Тест удаления пользователя
func TestUserHandler_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userID := uint(1)

	mockRepo.On("DeleteUser", userID).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/users/1", nil)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := handlers.NewUserHandler(mockRepo)
	handler.DeleteUser(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockRepo.AssertExpectations(t)
}
