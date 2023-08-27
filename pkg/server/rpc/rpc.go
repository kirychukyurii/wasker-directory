package rpc

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/kirychukyurii/wasker-directory/internal/controller"
	"github.com/kirychukyurii/wasker-directory/pkg/server/rpc/interceptor/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/kirychukyurii/wasker-directory/internal/config"
	"github.com/kirychukyurii/wasker-directory/pkg/logger"
	"github.com/kirychukyurii/wasker-directory/pkg/server/rpc/interceptor"
)

type Server struct {
	*grpc.Server
}

func New(cfg config.Config, log logger.Logger, controller controller.Controllers) Server {
	l, opts := interceptor.NewGrpcLoggingHandler(log)
	r := recovery.WithRecoveryHandler(interceptor.NewGrpcPanicRecoveryHandler(log))

	// create new gRPC server
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.ErrorUnaryServerInterceptor(),
			interceptor.RequestIdUnaryServerInterceptor(log),
			logging.UnaryServerInterceptor(l, opts...),
			auth.UnaryServerInterceptor(log, controller),
			recovery.UnaryServerInterceptor(r),
			// Add any other interceptor you want.
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(l, opts...),
			recovery.StreamServerInterceptor(r),
			// Add any other interceptor you want.
		))

	// Register reflection service on gRPC server.
	reflection.Register(s)

	return Server{s}
}
