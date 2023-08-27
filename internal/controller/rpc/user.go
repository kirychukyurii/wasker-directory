package rpc

import (
	"context"
	"time"

	v1 "buf.build/gen/go/kirychuk/wasker-proto/grpc/go/directory/v1/directoryv1grpc"
	common "buf.build/gen/go/kirychuk/wasker-proto/protocolbuffers/go/common/v1"
	pb "buf.build/gen/go/kirychuk/wasker-proto/protocolbuffers/go/directory/v1"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kirychukyurii/wasker-directory/internal/domain/entity"
	"github.com/kirychukyurii/wasker-directory/internal/domain/service"
	"github.com/kirychukyurii/wasker-directory/pkg/logger"
	"github.com/kirychukyurii/wasker-directory/pkg/model"
)

type UserService interface {
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
	ReadUser(ctx context.Context, userId int64) (*entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, userId int64) error
	QueryUsers(ctx context.Context, param *entity.UserQueryParam) (*entity.UserQueryResult, error)
}

// UserController will implement the service defined in protocol buffer definitions
type UserController struct {
	v1.UnimplementedUserServiceServer

	userService UserService
	log         logger.Logger
}

func NewUserController(userService *service.UserService, log logger.Logger) UserController {
	return UserController{
		userService: userService,
		log:         log,
	}
}

func (u UserController) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := entity.User{
		Model: model.Model{
			CreatedBy: model.LookupEntity{},
			UpdatedBy: model.LookupEntity{},
		},
		Name:     request.User.Name,
		Email:    request.User.Email,
		UserName: request.User.Username,
		Password: "",
		Role: &model.LookupEntity{
			Id: request.User.Role.Id,
		},
	}

	createdUser, err := u.userService.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	response := &pb.CreateUserResponse{
		User: &pb.User{
			Id:       createdUser.Id,
			Name:     createdUser.Name,
			Email:    createdUser.Email,
			Username: createdUser.UserName,
			Password: "",
			Role: &common.ObjectId{
				Id:   createdUser.Role.Id,
				Name: createdUser.Role.Name,
			},
			CreatedAt: createdUser.CreatedAt.Unix(),
			CreatedBy: &common.ObjectId{
				Id:   createdUser.CreatedBy.Id,
				Name: createdUser.CreatedBy.Name,
			},
			UpdatedAt: createdUser.UpdatedAt.Unix(),
			UpdatedBy: &common.ObjectId{
				Id:   createdUser.UpdatedBy.Id,
				Name: createdUser.UpdatedBy.Name,
			},
		},
	}

	return response, nil
}

func (u UserController) ReadUser(ctx context.Context, request *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	user, err := u.userService.ReadUser(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	response := &pb.ReadUserResponse{
		User: &pb.User{
			Id:       user.Id,
			Name:     user.Name,
			Email:    user.Email,
			Username: user.UserName,
			Password: "",
			Role: &common.ObjectId{
				Id:   user.Role.Id,
				Name: user.Role.Name,
			},
			CreatedAt: user.CreatedAt.Unix(),
			CreatedBy: &common.ObjectId{
				Id:   user.CreatedBy.Id,
				Name: user.CreatedBy.Name,
			},
			UpdatedAt: user.UpdatedAt.Unix(),
			UpdatedBy: &common.ObjectId{
				Id:   user.UpdatedBy.Id,
				Name: user.UpdatedBy.Name,
			},
		},
	}

	return response, nil
}

func (u UserController) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := entity.User{
		Model: model.Model{
			CreatedBy: model.LookupEntity{},
			UpdatedBy: model.LookupEntity{},
		},
		Name:     request.User.Name,
		Email:    request.User.Email,
		UserName: request.User.Username,
		Password: "",
		Role: &model.LookupEntity{
			Id: request.User.Role.Id,
		},
	}

	updatedUser, err := u.userService.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	response := &pb.UpdateUserResponse{
		Updated: &pb.User{
			Id:       updatedUser.Id,
			Name:     updatedUser.Name,
			Email:    updatedUser.Email,
			Username: updatedUser.UserName,
			Password: "",
			Role: &common.ObjectId{
				Id:   updatedUser.Role.Id,
				Name: updatedUser.Role.Name,
			},
			CreatedAt: updatedUser.CreatedAt.Unix(),
			CreatedBy: &common.ObjectId{
				Id:   updatedUser.CreatedBy.Id,
				Name: updatedUser.CreatedBy.Name,
			},
			UpdatedAt: updatedUser.UpdatedAt.Unix(),
			UpdatedBy: &common.ObjectId{
				Id:   updatedUser.UpdatedBy.Id,
				Name: updatedUser.UpdatedBy.Name,
			},
		},
	}

	return response, nil
}

func (u UserController) DeleteUsers(ctx context.Context, request *pb.DeleteUsersRequest) (*emptypb.Empty, error) {
	if err := u.userService.DeleteUser(ctx, request.GetId()); err != nil {
		return nil, err
	}

	return nil, nil
}

func (u UserController) ListUsers(ctx context.Context, request *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	param := &entity.UserQueryParam{
		Pagination: model.PaginationParam{
			Current:  request.Pagination.GetCurrent(),
			PageSize: request.Pagination.GetPageSize(),
		},
		Order: model.OrderParam{
			Key:       request.Order.GetKey(),
			Direction: model.OrderDirection(request.Order.GetDirection()),
		},
		Query: model.QueryParam{
			Id:        request.GetId(),
			Name:      request.GetName(),
			CreatedAt: time.Time{},
			CreatedBy: model.LookupEntity{},
			Query:     request.GetQ(),
		},
		UserName: request.Username,
	}

	users, err := u.userService.QueryUsers(ctx, param)
	if err != nil {
		return nil, err
	}

	list := users.List
	responseUsers := make([]*pb.User, len(list))
	for _, v := range list {
		user := &pb.User{
			Id:       v.Id,
			Name:     v.Name,
			Email:    v.Email,
			Username: v.UserName,
			Password: "",
			Role: &common.ObjectId{
				Id:   v.Role.Id,
				Name: v.Role.Name,
			},
			CreatedAt: v.CreatedAt.Unix(),
			CreatedBy: &common.ObjectId{
				Id:   v.CreatedBy.Id,
				Name: v.CreatedBy.Name,
			},
			UpdatedAt: v.UpdatedAt.Unix(),
			UpdatedBy: &common.ObjectId{
				Id:   v.UpdatedBy.Id,
				Name: v.UpdatedBy.Name,
			},
		}

		responseUsers = append(responseUsers, user)
	}

	response := &pb.ListUsersResponse{
		Items: responseUsers,
		Pagination: &common.Pagination{
			Next:     users.Pagination.Next,
			Current:  users.Pagination.Current,
			PageSize: users.Pagination.PageSize,
		},
	}

	return response, nil
}
