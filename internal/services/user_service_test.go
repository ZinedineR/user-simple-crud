package service_test

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"user-simple-crud/internal/entity"
	"user-simple-crud/internal/mocks"
	"user-simple-crud/internal/model"
	service "user-simple-crud/internal/services"
	mocksSignature "user-simple-crud/pkg/mocks"
	"user-simple-crud/pkg/xvalidator"
)

func setupSQLMock(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	// Setup SQL mock
	db, mockSql, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// Setup GORM with the mock DB
	gormDB, gormDBErr := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if gormDBErr != nil {
		t.Fatalf("failed to open GORM connection: %v", gormDBErr)
	}
	return mockSql, gormDB
}

func TestCreateUser(t *testing.T) {
	mockAppCtx := context.Background()

	t.Run("CreateUser Success", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
			Email:    "john@example.com",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(nil, nil)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "email", request.Email).Return(nil, nil)
		mockRepository.On("CreateTx", mockAppCtx, mock.Anything, mock.Anything).Return(nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("HashBscryptPassword", request.Password).Return("$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu", nil)

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		errService := mockService.Create(mockAppCtx, request)

		// Assert the result
		assert.Nil(t, errService)
	})

	t.Run("CreateUser Username and Email Empty", func(t *testing.T) {
		// Set up input (missing both username and email)
		request := &entity.UserLogin{
			Password: "SecurePass123!",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Create(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
	})

	t.Run("CreateUser Email Exists", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "SecurePass123!",
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:    "123e4567-e89b-12d3-a456-426614174000",
			Email: "john@example.com",
		}
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(nil, nil)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "email", request.Email).Return(existingUser, nil)

		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Create(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
	})
}

func TestLoginUser(t *testing.T) {
	mockAppCtx := context.Background()

	t.Run("LoginUser Success", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe",
			Password: "SecurePass123!",
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:       "123e4567-e89b-12d3-a456-426614174000",
			Username: "john_doe",
			Password: "$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu", // Hashed password
		}
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(existingUser, nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("CheckBscryptPasswordHash", request.Password, existingUser.Password).Return(true)
		mockSignaturer.On("GenerateJWT", existingUser.Username).Return("jwt_token", nil)

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.Login(mockAppCtx, request)

		// Assert the result
		assert.Nil(t, errService)
		assert.NotNil(t, result)
	})

	t.Run("LoginUser By Email Success", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Email:    "john@example.com",
			Password: "SecurePass123!",
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:       "123e4567-e89b-12d3-a456-426614174000",
			Email:    "john@example.com",
			Password: "$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu", // Hashed password
		}
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", "").Return(nil, nil)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "email", request.Email).Return(existingUser, nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("CheckBscryptPasswordHash", request.Password, existingUser.Password).Return(true)
		mockSignaturer.On("GenerateJWT", existingUser.Username).Return("jwt_token", nil)

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.Login(mockAppCtx, request)

		// Assert the result
		assert.Nil(t, errService)
		assert.NotNil(t, result)
		assert.Equal(t, "jwt_token", result.Token)
	})

	t.Run("LoginUser Username/Email Not Found", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "non_existent",
			Password: "SecurePass123!",
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(nil, nil)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "email", request.Email).Return(nil, nil)

		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.Login(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})
}

func TestUpdateUser(t *testing.T) {
	mockAppCtx := context.Background()

	t.Run("UpdateUser Success", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe_updated",
			Email:    "john_doe@example.com",
			Password: "NewSecurePass123!",
		}
		id := "123e4567-e89b-12d3-a456-426614174000"

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(nil, nil)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "email", request.Email).Return(nil, nil)
		mockRepository.On("UpdateTx", mockAppCtx, mock.Anything, mock.Anything).Return(nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("HashBscryptPassword", request.Password).Return("$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu", nil)

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		errService := mockService.Update(mockAppCtx, id, request)

		// Assert the result
		assert.Nil(t, errService)
	})

	t.Run("UpdateUser Invalid UUID", func(t *testing.T) {
		// Set up input with invalid UUID
		request := &entity.UserLogin{
			Username: "john_doe_updated",
			Email:    "john_doe@example.com",
			Password: "NewSecurePass123!",
		}
		id := "invalid-uuid"

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockSignaturer := new(mocksSignature.Signaturer)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Update(mockAppCtx, id, request)

		// Assert the result
		assert.NotNil(t, errService)
	})

	t.Run("UpdateUser Username Exists", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe_updated",
			Email:    "john_doe@example.com",
			Password: "NewSecurePass123!",
		}
		id := "123e4567-e89b-12d3-a456-426614174000"

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:       "different-id",
			Username: "john_doe_updated",
		}
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(existingUser, nil)
		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Update(mockAppCtx, id, request)

		// Assert the result
		assert.NotNil(t, errService)
	})

	t.Run("UpdateUser HashPassword Failed", func(t *testing.T) {
		// Set up input
		request := &entity.UserLogin{
			Username: "john_doe_updated",
			Email:    "john_doe@example.com",
			Password: "NewSecurePass123!",
		}
		id := "123e4567-e89b-12d3-a456-426614174000"

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "username", request.Username).Return(nil, nil)
		mockRepository.On("FindByName", mockAppCtx, mock.Anything, "email", request.Email).Return(nil, nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("HashBscryptPassword", request.Password).Return("", errors.New("hash error"))

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Update(mockAppCtx, id, request)

		// Assert the result
		assert.NotNil(t, errService)
	})
}

func TestDeleteUser(t *testing.T) {
	mockAppCtx := context.Background()

	t.Run("DeleteUser Success", func(t *testing.T) {
		// Set up input
		id := "123e4567-e89b-12d3-a456-426614174000"

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("DeleteByIDTx", mockAppCtx, mock.Anything, id).Return(nil)

		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		errService := mockService.Delete(mockAppCtx, id)

		// Assert the result
		assert.Nil(t, errService)
	})

	t.Run("DeleteUser Invalid UUID", func(t *testing.T) {
		// Set up input with invalid UUID
		id := "invalid-uuid"

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Delete(mockAppCtx, id)

		// Assert the result
		assert.NotNil(t, errService)
	})

	t.Run("DeleteUser Repository Error", func(t *testing.T) {
		// Set up input
		id := "123e4567-e89b-12d3-a456-426614174000"

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("DeleteByIDTx", mockAppCtx, mock.Anything, id).Return(errors.New("test error"))

		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		errService := mockService.Delete(mockAppCtx, id)

		// Assert the result
		assert.NotNil(t, errService)
	})
}

func TestFindOneUser(t *testing.T) {
	mockAppCtx := context.Background()

	t.Run("FindOneUser Success", func(t *testing.T) {
		// Set up input
		id := "123e4567-e89b-12d3-a456-426614174000"

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:       id,
			Username: "john_doe",
			Email:    "john_doe@example.com",
		}
		mockRepository.On("FindByID", mockAppCtx, mock.Anything, id).Return(existingUser, nil)
		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.FindOne(mockAppCtx, id)

		// Assert the result
		assert.Nil(t, errService)
		assert.NotNil(t, result)
	})

	t.Run("FindOneUser Invalid UUID", func(t *testing.T) {
		// Set up input with invalid UUID
		id := "invalid-uuid"

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.FindOne(mockAppCtx, id)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})

	t.Run("FindOneUser Repository Error", func(t *testing.T) {
		// Set up input
		id := "123e4567-e89b-12d3-a456-426614174000"

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByID", mockAppCtx, mock.Anything, id).Return(nil, errors.New("test error"))

		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.FindOne(mockAppCtx, id)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})
}

func TestListUser(t *testing.T) {
	mockAppCtx := context.Background()
	req := model.ListReq{
		Page: model.PaginationParam{
			Page:     1,
			PageSize: 1,
		},
		Order: model.OrderParam{
			Order:   "username",
			OrderBy: "asc",
		},
	}
	users := []*entity.User{
		{
			Id:       "0b8d3f3d-d343-4390-964c-4f05c4c803d6",
			Username: "john_doe",
			Email:    "john@example.com",
		},
	}

	t.Run("ListUser Success", func(t *testing.T) {
		// Setup the expected response from the repository
		response := &model.PaginationData[entity.User]{
			Page:             1,
			PageSize:         1,
			TotalPage:        1,
			TotalDataPerPage: 1,
			TotalData:        1,
			Data:             users,
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByPagination", mockAppCtx, mock.Anything, req.Page, req.Order, req.Filter).Return(response, nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.List(mockAppCtx, req)

		// Assert the result
		assert.Nil(t, errService)
		assert.NotNil(t, result)
	})

	t.Run("ListUser Failed Repository", func(t *testing.T) {
		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByPagination", mockAppCtx, mock.Anything, req.Page, req.Order, req.Filter).Return(nil, errors.New("test error"))
		mockSignaturer := new(mocksSignature.Signaturer)
		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.List(mockAppCtx, req)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})
}
