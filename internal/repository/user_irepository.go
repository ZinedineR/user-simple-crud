package repository

import (
	"context"
	"gorm.io/gorm"
	"user-simple-crud/internal/entity"
	"user-simple-crud/internal/model"
)

type UserRepository interface {
	// Example operations
	CreateTx(ctx context.Context, tx *gorm.DB, data *entity.User) error
	UpdateTx(ctx context.Context, tx *gorm.DB, data *entity.User) error
	FindByName(ctx context.Context, tx *gorm.DB, column, value string) (
		*entity.User, error,
	)
	FindByPagination(
		ctx context.Context, tx *gorm.DB, page model.PaginationParam, order model.OrderParam,
		filter model.FilterParams,
	) (*model.PaginationData[entity.User], error)
	FindByID(ctx context.Context, tx *gorm.DB, id string) (*entity.User, error)
	DeleteByIDTx(ctx context.Context, tx *gorm.DB, id string) error
}
