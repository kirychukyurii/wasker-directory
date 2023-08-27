package auth

import (
	"context"
	"fmt"
	"github.com/kirychukyurii/wasker-directory/internal/controller"
	"github.com/kirychukyurii/wasker-directory/pkg/server/rpc/interceptor"
	"github.com/kirychukyurii/wasker-directory/pkg/utils"
	"strings"

	"buf.build/gen/go/kirychuk/wasker-proto/grpc/go/directory/v1/directoryv1grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/kirychukyurii/wasker-directory/pkg/logger"
	"github.com/kirychukyurii/wasker-directory/pkg/werror"
)

var skipAuthServices = []string{
	directoryv1grpc.AuthService_ServiceDesc.ServiceName,
}

var (
	headerAuthorize = "authorization"
	typeAuthorize   = "bearer"
)

// UnaryServerInterceptor returns a server interceptor function to authenticate && authorize unary RPC
func UnaryServerInterceptor(log logger.Logger, controller controller.Controllers) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if ok := skipAuthInterceptor(info.FullMethod); !ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return "", werror.ErrRequestMissingMetadata
		}

		token, err := authFromMD(ctx, md[headerAuthorize], typeAuthorize)
		if err != nil {
			return nil, err
		}
		fmt.Println(token)

		/*
			userId, err := controller.Auth.Authn(ctx, token)
			if err != nil {
				return nil, err
			}

			ctx = context.WithValue(ctx, logger.LoggerCtxKey{}, log.FromContext(ctx).With().Int64("user", userId).Logger())

			service, method := splitFullMethodName(info.FullMethod)
			ok, err = controller.Auth.Authz(ctx, userId, service, method)
			if err != nil || !ok {
				return nil, err
			}
		*/

		return handler(ctx, req)
	}
}

// authFromMD is a helper function for extracting the :authorization header from the gRPC metadata of the request.
//
// It expects the `:authorization` header to be of a certain scheme (e.g. `basic`, `bearer`), in a
// case-insensitive format (see rfc2617, sec 1.2). If no such authorization is found, or the token
// is of wrong scheme, an error with gRPC status `Unauthenticated` is returned.
func authFromMD(ctx context.Context, md []string, expectedScheme string) (string, error) {
	if len(md) < 1 {
		return "", werror.NewUnauthenticatedError(werror.AppError{
			Message: werror.ErrAuthAccessTokenIncorrect.Error(),
			Details: werror.AppErrorDetail{
				Err:       werror.ErrAuthAccessTokenIncorrect,
				ErrReason: "NULL_METADATA",
				ErrDomain: "interceptor.auth.from_metadata",
				RequestId: utils.FromContext(ctx, interceptor.XRequestIDCtxKey{}).(string),
			},
		})
	}

	if md[0] == "" {
		return "", werror.NewUnauthenticatedError(werror.AppError{
			Message: werror.ErrAuthAccessTokenIncorrect.Error(),
			Details: werror.AppErrorDetail{
				Err:       werror.ErrAuthAccessTokenIncorrect,
				ErrReason: "EMPTY_METADATA",
				ErrDomain: "interceptor.auth.from_metadata",
				RequestId: utils.FromContext(ctx, interceptor.XRequestIDCtxKey{}).(string),
			},
		})
	}

	scheme, token, found := strings.Cut(md[0], " ")
	if !found {
		return "", werror.NewUnauthenticatedError(werror.AppError{
			Message: werror.ErrAuthAccessTokenIncorrect.Error(),
			Details: werror.AppErrorDetail{
				Err:       werror.ErrAuthAccessTokenIncorrect,
				ErrReason: "NULL_AUTHORIZATION_HEADER",
				ErrDomain: "interceptor.auth.from_metadata",
				RequestId: utils.FromContext(ctx, interceptor.XRequestIDCtxKey{}).(string),
			},
		})
	}

	if !strings.EqualFold(scheme, expectedScheme) {
		return "", werror.NewUnauthenticatedError(werror.AppError{
			Message: werror.ErrAuthAccessTokenIncorrect.Error(),
			Details: werror.AppErrorDetail{
				Err:       werror.ErrAuthAccessTokenIncorrect,
				ErrReason: "INVALID_TOKEN_TYPE",
				ErrDomain: "interceptor.auth.from_metadata",
				RequestId: utils.FromContext(ctx, interceptor.XRequestIDCtxKey{}).(string),
			},
		})
	}

	return token, nil
}

// skipAuthInterceptor setup auth matcher.
func skipAuthInterceptor(service string) bool {
	for _, s := range skipAuthServices {
		return s != service
	}

	return true
}

func splitFullMethodName(fullMethod string) (string, string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/") // remove leading slash
	if i := strings.Index(fullMethod, "/"); i >= 0 {
		return fullMethod[:i], fullMethod[i+1:]
	}

	return "unknown", "unknown"
}
