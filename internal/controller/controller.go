package controller

import (
	"github.com/kirychukyurii/wasker-directory/internal/controller/rpc"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewGroupControllers),
)

type Controllers struct {
	User rpc.UserController
	//Auth rpc.AuthController
}

func NewGroupControllers(u rpc.UserController /* a rpc.AuthController */) Controllers {
	return Controllers{
		User: u,
		//Auth: a,
	}
}
