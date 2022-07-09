package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"task/internal/app/grpc/api"
	"task/internal/app/services/configmanager"
)

type GRPCValidatorClient struct {
	config *configmanager.Config
}

func NewGRPCValidatorClient(config *configmanager.Config) *GRPCValidatorClient {
	return &GRPCValidatorClient{
		config: config,
	}
}

func (c *GRPCValidatorClient) Validate(accessToken string, refreshToken string) (*api.ValidateResponse, error) {
	conn, err := grpc.Dial(c.config.GRPCBindAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	cl := api.NewValidatorClient(conn)
	result, err := cl.Validate(context.Background(), &api.ValidateRequest{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	return result, err
}
