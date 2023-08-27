package interceptor

import (
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kirychukyurii/wasker-directory/pkg/logger"
)

func NewGrpcPanicRecoveryHandler(log logger.Logger) func(any) error {
	return grpcPanicRecoveryHandler(log)
}

func grpcPanicRecoveryHandler(log logger.Logger) func(any) error {
	return func(p any) (err error) {
		log.Error().Err(err).Msgf("recovered from panic: %s", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}
}
