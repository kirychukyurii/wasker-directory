package interceptor

import (
	"context"
	"fmt"
	"github.com/kirychukyurii/wasker-directory/pkg/utils"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"

	"github.com/kirychukyurii/wasker-directory/pkg/logger"
)

func NewGrpcLoggingHandler(log logger.Logger) (logging.Logger, []logging.Option) {
	return grpcLoggingHandler(log.Logger), grpcLoggingOption()
}

// GrpcLoggingHandler adapts zerolog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func grpcLoggingHandler(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l := l.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			l.Debug().Msg(msg)
		case logging.LevelInfo:
			l.Info().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func grpcLoggingOption() []logging.Option {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.FinishCall),
		logging.WithFieldsFromContext(func(ctx context.Context) logging.Fields {
			return logging.Fields{XRequestIDKey, utils.FromContext(ctx, XRequestIDCtxKey{})}
		}),
		// Add any other option (check functions starting with logging.With).
	}

	return opts
}
