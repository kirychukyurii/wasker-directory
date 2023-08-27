package rpc

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserController),
	//fx.Provide(NewAuthController),
)
