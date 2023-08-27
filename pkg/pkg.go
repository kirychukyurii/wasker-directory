package pkg

import (
	"go.uber.org/fx"

	"github.com/kirychukyurii/wasker-directory/pkg/consul"
	"github.com/kirychukyurii/wasker-directory/pkg/db"
	"github.com/kirychukyurii/wasker-directory/pkg/logger"
	"github.com/kirychukyurii/wasker-directory/pkg/server/rpc"
)

// Module exports dependency
var Module = fx.Options(
	fx.Provide(logger.New),
	fx.Provide(db.New),
	fx.Provide(consul.New),
	fx.Provide(rpc.New),
)
