package Server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/DeebTibi/GoVault/config"
	RegistryClient "github.com/DeebTibi/GoVault/services/registry/client"
	service "github.com/DeebTibi/GoVault/services/user_auth/api"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type userAuthServer struct {
	service.UserAuthServer
}

func NewUserAuthServer() service.UserAuthServer {
	return &userAuthServer{}
}

func Start(cfg *config.ServiceConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", "localhost", 0))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.ConnectionTimeout(time.Second*10))
	grpcServer := grpc.NewServer(opts...)
	service.RegisterUserAuthServer(grpcServer, NewUserAuthServer())
	registryClient := RegistryClient.NewRegistryClient()
	registryClient.Register("user_auth", lis.Addr().String())
	defer registryClient.Unregister("user_auth", lis.Addr().String())
	grpcServer.Serve(lis)
}

func (obj *userAuthServer) Login(ctx context.Context, in *service.LoginRequest) (*service.LoginResponse, error) {
	res, err := LoginUser(in.UserName, in.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}
	return &service.LoginResponse{Token: res}, nil
}

func (obj *userAuthServer) Register(ctx context.Context, in *service.RegisterRequest) (*service.RegisterResponse, error) {
	res, err := RegisterUser(in.UserName, in.Password, in.UserKey)
	if err != nil {
		return nil, err
	}
	return &service.RegisterResponse{Token: res}, nil
}

func (obj *userAuthServer) AuthenticateToken(ctx context.Context, in *service.AuthenticateTokenRequest) (*wrappers.BoolValue, error) {
	res, err := AuthenticateUser(in.UserName, in.Token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}
	return &wrapperspb.BoolValue{Value: res}, nil
}
