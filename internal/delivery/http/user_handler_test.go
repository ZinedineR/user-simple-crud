package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-simple-crud/internal/entity"
	"user-simple-crud/internal/mocks"
	service "user-simple-crud/internal/services"
	"user-simple-crud/pkg/exception"
)

func TestUserHttpHandler_Register(t *testing.T) {
	t.Run("Create Success", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/auth/register", userHandler.Register)

		// Mock Data
		requestBody := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service call
		mockUserService.On("Create", mock.Anything, requestBody).Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Create Error - Invalid JSON", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/auth/register", userHandler.Register)

		// Malformed JSON
		malformedJSON := `{"invalid_json"}`
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(malformedJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Create Error - Service Error", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/auth/register", userHandler.Register)

		// Mock Data
		requestBody := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service call with error
		mockUserService.On("Create", mock.Anything, requestBody).Return(exception.Internal("error", errors.New("registration failed")))

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserHttpHandler_Login(t *testing.T) {
	t.Run("Login Success", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/auth/login", userHandler.Login)

		// Mock Data
		requestBody := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		expectedResponse := &service.UserLoginResponse{
			Username: "john_doe",
			Token:    "jwt_token",
		}

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service call
		mockUserService.On("Login", mock.Anything, requestBody).Return(expectedResponse, nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Login Error - Invalid JSON", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/auth/login", userHandler.Login)

		// Malformed JSON
		malformedJSON := `{"invalid_json"}`
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(malformedJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Login Error - Service Error", func(t *testing.T) {
		// Setup
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/auth/login", userHandler.Login)

		// Mock Data
		requestBody := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock service call with error
		mockUserService.On("Login", mock.Anything, requestBody).Return(nil, exception.Internal("error", errors.New("login failed")))

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserHttpHandler_Create(t *testing.T) {
	// Setup router
	t.Run("CreateUser Success", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/users", userHandler.Create)

		// Prepare request data
		requestBody := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUserService.On("Create", mock.Anything, requestBody).Return(nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("CreateUser Service Error", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/users", userHandler.Create)

		requestBody := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Create HTTP POST request
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Set up the expectation on the mock service
		mockUserService.On("Create", mock.Anything, requestBody).Return(exception.Internal("error", errors.New("test error")))

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("CreateUser Binding JSON Error", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.POST("/users", userHandler.Create)

		// Malformed JSON
		malformedJSON := `{"invalid_json"`
		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(malformedJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserHttpHandler_List(t *testing.T) {
	expectResponse := &service.ListUserResp{}

	t.Run("ListUsers Success", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.GET("/users", userHandler.List)

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/users", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock the service
		mockUserService.On("List", mock.Anything, mock.Anything).Return(expectResponse, nil)

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ListUsers Service Error", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.GET("/users", userHandler.List)

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/users", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Mock the service
		mockUserService.On("List", mock.Anything, mock.Anything).Return(nil, exception.Internal("error", errors.New("test error")))

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("ListUsers BadRequest Error", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.GET("/users", userHandler.List)

		// Simulate a bad request with an invalid filter
		req, _ := http.NewRequest("GET", "/users?filter=id:invalid:error", nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserHttpHandler_FindOne(t *testing.T) {
	t.Run("FindOneUser Success", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.GET("/users/:id", userHandler.FindOne)

		// Mock Data
		userID := "123e4567-e89b-12d3-a456-426614174000"
		expectedUser := &entity.User{
			Id:       userID,
			Username: "john_doe",
			Email:    "john_doe@example.com",
		}

		// Mock the service
		mockUserService.On("FindOne", mock.Anything, userID).Return(expectedUser, nil)

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/users/"+userID, nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("FindOneUser Not Found", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.GET("/users/:id", userHandler.FindOne)

		// Mock Data
		userID := "invalid-id"
		mockUserService.On("FindOne", mock.Anything, userID).Return(nil, exception.NotFound("user not found"))

		// Create HTTP GET request
		req, _ := http.NewRequest("GET", "/users/"+userID, nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUserHttpHandler_Update(t *testing.T) {
	t.Run("UpdateUser Success", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.PUT("/users/:id", userHandler.Update)

		// Mock Data
		userID := "123e4567-e89b-12d3-a456-426614174000"
		requestBody := &entity.UserLogin{
			Username: "john_doe_updated",
			Password: "NewSecurePass123!",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Mock the service
		mockUserService.On("Update", mock.Anything, userID, requestBody).Return(nil)

		// Create HTTP PUT request
		req, _ := http.NewRequest("PUT", "/users/"+userID, bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UpdateUser Binding JSON Error", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.PUT("/users/:id", userHandler.Update)

		// Malformed JSON
		malformedJSON := `{"invalid_json"}`
		userID := "123e4567-e89b-12d3-a456-426614174000"

		// Create HTTP PUT request
		req, _ := http.NewRequest("PUT", "/users/"+userID, bytes.NewBufferString(malformedJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("UpdateUser Service Error", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.PUT("/users/:id", userHandler.Update)

		// Mock Data
		userID := "123e4567-e89b-12d3-a456-426614174000"
		requestBody := &entity.UserLogin{
			Username: "john_doe_updated",
			Password: "NewSecurePass123!",
		}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Mock the service
		mockUserService.On("Update", mock.Anything, userID, requestBody).Return(exception.Internal("error", errors.New("update failed")))

		// Create HTTP PUT request
		req, _ := http.NewRequest("PUT", "/users/"+userID, bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserHttpHandler_Delete(t *testing.T) {
	t.Run("DeleteUser Success", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.DELETE("/users/:id", userHandler.Delete)

		// Mock Data
		userID := "123e4567-e89b-12d3-a456-426614174000"

		// Mock the service
		mockUserService.On("Delete", mock.Anything, userID).Return(nil)

		// Create HTTP DELETE request
		req, _ := http.NewRequest("DELETE", "/users/"+userID, nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DeleteUser Service Error", func(t *testing.T) {
		r := gin.Default()
		mockUserService := new(mocks.UserService)
		userHandler := NewUserHTTPHandler(mockUserService)

		r.DELETE("/users/:id", userHandler.Delete)

		// Mock Data
		userID := "123e4567-e89b-12d3-a456-426614174000"

		// Mock the service
		mockUserService.On("Delete", mock.Anything, userID).Return(exception.Internal("error", errors.New("delete failed")))

		// Create HTTP DELETE request
		req, _ := http.NewRequest("DELETE", "/users/"+userID, nil)
		req.Header.Set("Content-Type", "application/json")

		// Create gin context
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = req

		// Perform request
		r.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
