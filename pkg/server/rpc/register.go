package rpc

import (
	pb "buf.build/gen/go/kirychuk/wasker-proto/grpc/go/directory/v1/directoryv1grpc"
	"google.golang.org/grpc"

	"github.com/kirychukyurii/wasker-directory/internal/controller"
)

func GrpcDirectoryServiceServers(s grpc.ServiceRegistrar, controller controller.Controllers) {
	pb.RegisterUserServiceServer(s, &controller.User)
	//pb.RegisterAuthServiceServer(s, &controller.Auth)
}
