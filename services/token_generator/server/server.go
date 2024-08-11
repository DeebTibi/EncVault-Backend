package Server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/DeebTibi/GoVault/config"
	RegistryClient "github.com/DeebTibi/GoVault/services/registry/client"
	service "github.com/DeebTibi/GoVault/services/token_generator/api"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type tokenGeneratorServerImplementation struct {
	service.TokenGeneratorServer
}

func NewTokenGeneratorServer() service.TokenGeneratorServer {
	return &tokenGeneratorServerImplementation{}
}

func Start(cfg *config.ServiceConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", "localhost", 0))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.ConnectionTimeout(time.Second*10))
	grpcServer := grpc.NewServer(opts...)
	service.RegisterTokenGeneratorServer(grpcServer, NewTokenGeneratorServer())
	registryClient := RegistryClient.NewRegistryClient()
	registryClient.Register("token_generator", lis.Addr().String())
	defer registryClient.Unregister("token_generator", lis.Addr().String())
	grpcServer.Serve(lis)
}

func (s *tokenGeneratorServerImplementation) CreateUserToken(ctx context.Context, in *service.CreateUserTokenRequest) (*wrappers.StringValue, error) {
	token, err := CreateUserToken(in.UserId)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{Value: token}, nil
}

func (s *tokenGeneratorServerImplementation) ValidateUserToken(ctx context.Context, in *service.ValidateUserTokenRequest) (*wrappers.BoolValue, error) {
	valid := ValidateUserToken(in.UserId, in.UserToken)
	return &wrapperspb.BoolValue{Value: valid}, nil
}
