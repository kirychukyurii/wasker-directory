package auth

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"buf.build/gen/go/kirychuk/wasker-proto/grpc/go/directory/v1/directoryv1grpc"
)

type Client struct {
	Auth directoryv1grpc.AuthServiceClient
}

func New() (*Client, error) {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := directoryv1grpc.NewAuthServiceClient(conn)

	return &Client{
		Auth: client,
	}, nil
}
