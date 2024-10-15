package repository

import (
	"user-simple-crud/internal/entity"
)

type UserSQLRepo struct {
	Repository[entity.User]
}

func NewUserSQLRepository() UserRepository {
	return &UserSQLRepo{}
}
