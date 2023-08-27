package interceptor

import (
	"context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"

	"github.com/kirychukyurii/wasker-directory/pkg/werror"
	"google.golang.org/grpc"
)

// ErrorUnaryServerInterceptor returns a server interceptor function to authenticate && authorize unary RPC
func ErrorUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		h, err := handler(ctx, req)
		if err != nil {
			appErr, ok := err.(*werror.AppError)
			if ok {
				s := status.New(appErr.Code, appErr.Message)
				sd, wdErr := s.WithDetails(&errdetails.ErrorInfo{
					Reason: appErr.Details.ErrReason,
					Domain: appErr.Details.ErrDomain,
					Metadata: map[string]string{
						"request_id": appErr.Details.RequestId,
					},
				})

				if wdErr != nil {
					return h, s.Err()
				}

				return h, sd.Err()
			} else {
				return h, err
			}
		}
		return h, nil
	}
}
