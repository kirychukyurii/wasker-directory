package run

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/fx"

	"github.com/kirychukyurii/wasker-directory/internal/adapter/storage"
	"github.com/kirychukyurii/wasker-directory/internal/config"
	"github.com/kirychukyurii/wasker-directory/internal/controller"
	crpc "github.com/kirychukyurii/wasker-directory/internal/controller/rpc"
	"github.com/kirychukyurii/wasker-directory/internal/domain/service"
	"github.com/kirychukyurii/wasker-directory/pkg"
	"github.com/kirychukyurii/wasker-directory/pkg/consul"
	"github.com/kirychukyurii/wasker-directory/pkg/db"
	"github.com/kirychukyurii/wasker-directory/pkg/logger"
	"github.com/kirychukyurii/wasker-directory/pkg/server/rpc"
	"github.com/kirychukyurii/wasker-directory/pkg/uuid"
)

var Module = fx.Options(config.Module,
	pkg.Module,
	storage.Module,
	service.Module,
	crpc.Module,
	controller.Module,
	fx.Invoke(runApplication),
)

func runApplication(lifecycle fx.Lifecycle, cfg config.Config, log logger.Logger, db db.Database,
	grpcServer rpc.Server, discovery consul.ServiceDiscovery, controller controller.Controllers) {
	serviceId := fmt.Sprintf("%s-%s", ServiceName, uuid.New().String())
	subLogger := log.With().Str("service-id", serviceId).Logger()

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			subLogger.Info().Str("grpc-listen", cfg.Grpc.ListenAddr()).Msg("starting application")
			svcReg := consul.ServiceRegistration{
				Id:       serviceId,
				Service:  ServiceName,
				Host:     cfg.Grpc.Host,
				Port:     cfg.Grpc.Port,
				Protocol: "http2",
			}

			err := discovery.Register(svcReg)
			if err != nil {
				subLogger.Fatal().Err(err).Msg("register service")
			}

			go func() {
				l, err := net.Listen("tcp", cfg.Grpc.ListenAddr())
				if err != nil {
					subLogger.Fatal().Err(err).Msgf("listening on port :%d", cfg.Grpc.Port)
				}

				// register controllers on the gRPC server
				rpc.GrpcDirectoryServiceServers(grpcServer, controller)

				// the gRPC server
				if err := grpcServer.Serve(l); err != nil {
					subLogger.Fatal().Err(err).Msg("start server")
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			subLogger.Info().Msg("stopping application")

			db.Pool.Close()
			grpcServer.GracefulStop()
			err := discovery.Client.Agent().ServiceDeregister(serviceId)
			if err != nil {
				subLogger.Error().Err(err).Msg("deregister service")
			}

			return nil
		},
	})
}
