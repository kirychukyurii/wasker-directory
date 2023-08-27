package interceptor

import (
	"context"
	"github.com/kirychukyurii/wasker-directory/pkg/utils"
	"github.com/kirychukyurii/wasker-directory/pkg/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/kirychukyurii/wasker-directory/pkg/logger"
)

type XRequestIDCtxKey struct{}

// XRequestIDKey is metadata key name for request ID
var XRequestIDKey = "x-request-id"

func RequestIdUnaryServerInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		reqId := handleRequestID(ctx)

		childCtx := utils.ContextWithValue(ctx, XRequestIDCtxKey{}, reqId)
		childCtx = utils.ContextWithValue(childCtx, logger.LoggerCtxKey{}, log.With().Str(XRequestIDKey, reqId).Logger())

		return handler(childCtx, req)
	}
}

func handleRequestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.New().String()
	}

	header, ok := md[XRequestIDKey]
	if !ok || len(header) == 0 {
		return uuid.New().String()
	}

	requestID := header[0]
	if requestID == "" {
		return uuid.New().String()
	}

	return requestID
}
