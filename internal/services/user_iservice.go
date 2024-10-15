package service

import (
	"context"
	"user-simple-crud/internal/entity"
	"user-simple-crud/internal/model"
	"user-simple-crud/pkg/exception"
)

type UserService interface {
	// Register-Login operations for User
	Create(
		ctx context.Context, model *entity.UserLogin,
	) *exception.Exception
	Login(ctx context.Context, model *entity.UserLogin) (*UserLoginResponse, *exception.Exception)

	// CRUD operations for User
	Update(
		ctx context.Context, id string, model *entity.UserLogin,
	) *exception.Exception
	Delete(
		ctx context.Context, id string,
	) *exception.Exception
	List(ctx context.Context, req model.ListReq) (
		*ListUserResp, *exception.Exception,
	)
	FindOne(ctx context.Context, id string) (*entity.User, *exception.Exception)
}

type UserLoginResponse struct {
	Username string `json:"username" example:"john_doe"`
	Email    string `json:"email" example:"john_doe@example.com"`
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"` // JWT token example
}

type ListUserResp struct {
	Pagination *model.Pagination `json:"pagination"`
	Data       []*entity.User    `json:"data"`
}
