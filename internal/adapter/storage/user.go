package storage

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	scan "github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"github.com/kirychukyurii/wasker-directory/internal/domain/entity"
	"github.com/kirychukyurii/wasker-directory/pkg/db"
	"github.com/kirychukyurii/wasker-directory/pkg/logger"
	"github.com/kirychukyurii/wasker-directory/pkg/model"
	"github.com/kirychukyurii/wasker-directory/pkg/server/rpc/interceptor"
	"github.com/kirychukyurii/wasker-directory/pkg/utils"
	"github.com/kirychukyurii/wasker-directory/pkg/werror"
)

type UserStorage struct {
	db  db.Database
	log logger.Logger
}

func NewUserStorage(db db.Database, log logger.Logger) *UserStorage {
	return &UserStorage{
		db:  db,
		log: log,
	}
}

func (a UserStorage) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (a UserStorage) ReadUser(ctx context.Context, userId int64) (*entity.User, error) {
	var user entity.User

	q := a.db.Dialect().Select("u.id", "coalesce(u.name, u.user_name) AS name", "u.user_name", "u.password", "u.email",
		"u.created_at", `u.created_by "created_by.id"`, `coalesce(c.name, c.user_name) "created_by.name"`,
		"u.updated_at", `u.updated_by "updated_by.id"`, `coalesce(upd.name, upd.user_name) "updated_by.name"`,
		`u.role_id "role.id"`, `r.name "role.name"`).From("auth_user u").
		LeftJoin("auth_role r ON r.id = u.role_id").
		InnerJoin("auth_user c ON c.id = u.created_by").
		InnerJoin("auth_user upd ON upd.id = u.updated_by").
		Where(sq.Eq{"u.deleted_at": nil}).Where(sq.Eq{"u.id": userId})

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, werror.NewInternalError(werror.AppError{
			Message: werror.ErrDatabaseInternalError.Error(),
			Details: werror.AppErrorDetail{
				Err:       err,
				ErrReason: werror.ErrBuildQueryReason,
				ErrDomain: "repository.user.read",
				RequestId: utils.FromContext(ctx, interceptor.XRequestIDCtxKey{}).(string),
			},
		})
	}

	if err = scan.Get(ctx, a.db.Pool, &user, sql, args...); err != nil {
		var dbErr error
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			dbErr = werror.ErrDatabaseRecordNotFound
		default:
			dbErr = werror.ErrDatabaseInternalError
		}

		return nil, werror.NewInternalError(werror.AppError{
			Message: dbErr.Error(),
			Details: werror.AppErrorDetail{
				Err:       err,
				ErrReason: werror.ErrExecQueryReason,
				ErrDomain: "repository.user.read",
				RequestId: utils.FromContext(ctx, interceptor.XRequestIDCtxKey{}).(string),
			},
		})
	}

	return &user, nil
}

func (a UserStorage) UpdateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (a UserStorage) DeleteUser(ctx context.Context, userId int64) error {
	//TODO implement me
	panic("implement me")
}

func (a UserStorage) QueryUsers(ctx context.Context, param *entity.UserQueryParam) (*entity.UserQueryResult, error) {
	var list entity.Users
	var pagination model.Pagination

	q := a.db.Dialect().Select("u.id", "coalesce(u.name, u.user_name) AS name", "u.user_name", "u.password", "u.email",
		"u.created_at", `u.created_by "created_by.id"`, `coalesce(c.name, c.user_name) "created_by.name"`,
		"u.updated_at", `u.updated_by "updated_by.id"`, `coalesce(upd.name, upd.user_name) "updated_by.name"`,
		`u.role_id "role.id"`, `r.name "role.name"`).From("auth_user u").
		LeftJoin("auth_role r ON r.id = u.role_id").
		InnerJoin("auth_user c ON c.id = u.created_by").
		InnerJoin("auth_user upd ON upd.id = u.updated_by").
		Where(sq.Eq{"u.deleted_at": nil})

	if v := param.Query.Id; len(v) != 0 {
		q = q.Where(sq.Eq{"u.id": v})
	}

	if v := param.Query.Name; v != "" {
		q = q.Where(sq.Eq{"u.name": v})
	}

	if v := param.UserName; v != "" {
		q = q.Where(sq.Eq{"u.user_name": v})
	}

	q = q.OrderBy(param.Order.Parse())
	current, pageSize := param.Pagination.GetCurrent(), param.Pagination.GetPageSize()
	if current > 0 && pageSize > 0 {
		offset := (current - 1) * pageSize
		q = q.Offset(offset).Limit(pageSize)
	} else if pageSize > 0 {
		q = q.Limit(pageSize)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, werror.NewInternalError(werror.AppError{
			Message: werror.ErrDatabaseInternalError.Error(),
			Details: werror.AppErrorDetail{
				Err:       err,
				ErrReason: werror.ErrBuildQueryReason,
				ErrDomain: "repository.user.query",
				RequestId: utils.FromContext(ctx, interceptor.XRequestIDCtxKey{}).(string),
			},
		})
	}

	if err = scan.Select(ctx, a.db.Pool, &list, sql, args...); err != nil {
		var dbErr error
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			dbErr = werror.ErrDatabaseRecordNotFound
		default:
			dbErr = werror.ErrDatabaseInternalError
		}

		return nil, werror.NewInternalError(werror.AppError{
			Message: dbErr.Error(),
			Details: werror.AppErrorDetail{
				Err:       err,
				ErrReason: werror.ErrExecQueryReason,
				ErrDomain: "repository.user.query",
				RequestId: utils.FromContext(ctx, interceptor.XRequestIDCtxKey{}).(string),
			},
		})
	}

	pagination.Current = current
	pagination.PageSize = pageSize
	qr := &entity.UserQueryResult{
		Pagination: &pagination,
		List:       list,
	}

	return qr, nil
}
