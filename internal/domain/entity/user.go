package entity

import "github.com/kirychukyurii/wasker-directory/pkg/model"

type User struct {
	model.Model

	Name     string              `json:"name" storage:"name"`
	Email    string              `json:"email" validate:"required,email" storage:"email"`
	UserName string              `json:"username" validate:"required" storage:"user_name"`
	Password string              `json:"password" validate:"required" storage:"password"`
	Role     *model.LookupEntity `json:"role" storage:"role"`
}

type Users []*User

type (
	UserQueryParam struct {
		Pagination model.PaginationParam
		Order      model.OrderParam
		Query      model.QueryParam

		UserName string `query:"user_name"`
	}

	UserQueryResult struct {
		List       Users             `json:"list"`
		Pagination *model.Pagination `json:"pagination"`
	}
)
