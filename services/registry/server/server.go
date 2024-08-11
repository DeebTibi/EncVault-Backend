package Server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/DeebTibi/GoVault/config"
	service "github.com/DeebTibi/GoVault/services/registry/api"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type registryService struct {
	service.RegistryServer
}

func CreateNewRegistryService() service.RegistryServer {
	return &registryService{}
}

func Start(cfg *config.ServiceConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", "localhost", 8502))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.ConnectionTimeout(time.Second*10))
	grpcServer := grpc.NewServer(opts...)
	service.RegisterRegistryServer(grpcServer, CreateNewRegistryService())
	fmt.Printf("Starting registry service on %v\n", 8502)
	grpcServer.Serve(lis)
}

func (obj *registryService) Discover(ctx context.Context, in *service.DiscoverRequest) (*service.DiscoverResponse, error) {
	res := GetServiceIps(in.ServiceName)
	if res == nil {
		return nil, fmt.Errorf("service not found")
	}
	return &service.DiscoverResponse{ServiceIps: res}, nil
}

func (obj *registryService) Register(ctx context.Context, in *service.RegisterRequest) (*empty.Empty, error) {
	RegisterService(in.ServiceName, in.ServiceIp)
	return &empty.Empty{}, nil
}

func (obj *registryService) Unregister(ctx context.Context, in *service.UnregisterRequest) (*empty.Empty, error) {
	UnregisterService(in.ServiceName, in.ServiceIp)
	return &empty.Empty{}, nil
}
