package service

import (
	"context"
	"github.com/google/uuid"
	"user-simple-crud/internal/entity"
	"user-simple-crud/internal/model"
	"user-simple-crud/internal/repository"
	"user-simple-crud/pkg/signature"

	//"user-simple-crud/pkg/exception"
	"gorm.io/gorm"
	"user-simple-crud/pkg/exception"
	"user-simple-crud/pkg/xvalidator"
)

type UserServiceImpl struct {
	db         *gorm.DB
	userRepo   repository.UserRepository
	signaturer signature.Signaturer
	validate   *xvalidator.Validator
}

func NewUserService(
	db *gorm.DB, repo repository.UserRepository,
	signaturer signature.Signaturer,
	validate *xvalidator.Validator,
) UserService {
	return &UserServiceImpl{
		db:         db,
		userRepo:   repo,
		signaturer: signaturer,
		validate:   validate,
	}
}

func (s *UserServiceImpl) Create(
	ctx context.Context, model *entity.UserLogin,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(model); errs != nil {
		return exception.InvalidArgument(errs)
	}
	if model.Email == "" && model.Username == "" {
		return exception.InvalidArgument("either email or username must be filled")
	}
	duplicateCheck, err := s.userRepo.FindByName(ctx, s.db, "username", model.Username)
	if err != nil {
		return exception.Internal("err", err)
	}
	if duplicateCheck != nil {
		return exception.PermissionDenied("username already exists")
	}
	duplicateCheck, err = s.userRepo.FindByName(ctx, s.db, "email", model.Email)
	if err != nil {
		return exception.Internal("err", err)
	}
	if duplicateCheck != nil {
		return exception.PermissionDenied("email already exists")
	}
	password, err := s.signaturer.HashBscryptPassword(model.Password)
	if err != nil {
		return exception.Internal("can't create password", err)
	}
	body := &entity.User{
		Id:       uuid.NewString(),
		Username: model.Username,
		Email:    model.Email,
		Password: password,
	}
	if err := s.userRepo.CreateTx(ctx, tx, body); err != nil {
		return exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *UserServiceImpl) Login(ctx context.Context, model *entity.UserLogin) (
	*UserLoginResponse, *exception.Exception,
) {
	if errs := s.validate.Struct(model); errs != nil {
		return nil, exception.InvalidArgument(errs)
	}
	result, err := s.userRepo.FindByName(ctx, s.db, "username", model.Username)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if result == nil {
		result, err = s.userRepo.FindByName(ctx, s.db, "email", model.Email)
		if err != nil {
			return nil, exception.Internal("err", err)
		}
		if result == nil {
			return nil, exception.NotFound("username/email not found")
		}
	}
	if ok := s.signaturer.CheckBscryptPasswordHash(model.Password, result.Password); !ok {
		return nil, exception.PermissionDenied("username/password unmatched")
	}
	jwtToken, err := s.signaturer.GenerateJWT(result.Username)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	return &UserLoginResponse{
		Username: result.Username,
		Token:    jwtToken,
	}, nil
}

func (s *UserServiceImpl) Update(
	ctx context.Context, id string, model *entity.UserLogin,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(model); errs != nil {
		return exception.InvalidArgument(errs)
	}
	_, err := uuid.Parse(id)
	if err != nil {
		return exception.InvalidArgument("invalid user id, must be uuid")
	}
	if model.Email == "" && model.Username == "" {
		return exception.InvalidArgument("either email or username must be filled")
	}
	duplicateCheck, err := s.userRepo.FindByName(ctx, s.db, "username", model.Username)
	if err != nil {
		return exception.Internal("err", err)
	}
	if duplicateCheck != nil && duplicateCheck.Id != id {
		return exception.PermissionDenied("username already exists")
	}
	duplicateCheck, err = s.userRepo.FindByName(ctx, s.db, "email", model.Email)
	if err != nil {
		return exception.Internal("err", err)
	}
	if duplicateCheck != nil && duplicateCheck.Id != id {
		return exception.PermissionDenied("email already exists")
	}
	password, err := s.signaturer.HashBscryptPassword(model.Password)
	if err != nil {
		return exception.Internal("can't create password", err)
	}
	body := &entity.User{
		Id:       id,
		Username: model.Username,
		Email:    model.Email,
		Password: password,
	}
	if err := s.userRepo.UpdateTx(ctx, tx, body); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *UserServiceImpl) Delete(
	ctx context.Context, id string,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	_, err := uuid.Parse(id)
	if err != nil {
		return exception.InvalidArgument("invalid user id, must be uuid")
	}
	if err := s.userRepo.DeleteByIDTx(ctx, tx, id); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *UserServiceImpl) List(ctx context.Context, req model.ListReq) (
	*ListUserResp, *exception.Exception,
) {
	result, err := s.userRepo.FindByPagination(ctx, s.db, req.Page, req.Order, req.Filter)
	if err != nil {
		return nil, exception.Internal("failed to get User", err)
	}
	return &ListUserResp{
		Pagination: &model.Pagination{
			Page:             result.Page,
			PageSize:         result.PageSize,
			TotalPage:        result.TotalPage,
			TotalDataPerPage: result.TotalDataPerPage,
			TotalData:        result.TotalData,
		},
		Data: result.Data,
	}, nil
}

func (s *UserServiceImpl) FindOne(ctx context.Context, id string) (*entity.User, *exception.Exception) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, exception.InvalidArgument("invalid user id, must be uuid")
	}
	result, err := s.userRepo.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	return result, nil
}
