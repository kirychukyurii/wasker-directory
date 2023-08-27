package rpc

import (
	pb "buf.build/gen/go/kirychuk/wasker-proto/grpc/go/directory/v1/directoryv1grpc"
	"github.com/kirychukyurii/wasker-directory/pkg/logger"
)

type AuthService interface {
}

// AuthController will implement the service defined in protocol buffer definitions
type AuthController struct {
	pb.UnimplementedAuthServiceServer

	authService AuthService
	log         logger.Logger
}

func NewAuthController(authService AuthService, log logger.Logger) AuthController {
	return AuthController{
		authService: authService,
		log:         log,
	}
}
