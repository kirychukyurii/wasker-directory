package service

import (
	"context"
	"github.com/kirychukyurii/wasker-directory/internal/adapter/storage"

	"github.com/kirychukyurii/wasker-directory/internal/domain/entity"
	"github.com/kirychukyurii/wasker-directory/pkg/logger"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
	ReadUser(ctx context.Context, userId int64) (*entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, userId int64) error
	QueryUsers(ctx context.Context, param *entity.UserQueryParam) (*entity.UserQueryResult, error)
}

type UserService struct {
	log         logger.Logger
	userStorage UserStorage
}

func NewUserService(userStorage *storage.UserStorage, log logger.Logger) *UserService {
	return &UserService{
		log:         log,
		userStorage: userStorage,
	}
}

func (u UserService) ReadUser(ctx context.Context, userId int64) (*entity.User, error) {
	user, err := u.userStorage.ReadUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u UserService) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	createdUser, err := u.userStorage.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u UserService) UpdateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	updatedUser, err := u.userStorage.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u UserService) DeleteUser(ctx context.Context, userId int64) error {
	if err := u.userStorage.DeleteUser(ctx, userId); err != nil {
		return err
	}

	return nil
}

func (u UserService) QueryUsers(ctx context.Context, param *entity.UserQueryParam) (*entity.UserQueryResult, error) {
	users, err := u.userStorage.QueryUsers(ctx, param)
	if err != nil {
		return nil, err
	}

	return users, nil
}
