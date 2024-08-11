package Server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/DeebTibi/GoVault/config"
	service "github.com/DeebTibi/GoVault/services/key_maker/api"
	RegistryClient "github.com/DeebTibi/GoVault/services/registry/client"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type keyMakerServerImplementation struct {
	service.KeyMakerServer
}

func NewKeyMakerServer() service.KeyMakerServer {
	return &keyMakerServerImplementation{}
}

func Start(cfg *config.ServiceConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", "localhost", 0))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.ConnectionTimeout(time.Second*10))
	grpcServer := grpc.NewServer(opts...)
	service.RegisterKeyMakerServer(grpcServer, NewKeyMakerServer())
	registryClient := RegistryClient.NewRegistryClient()
	registryClient.Register("key_maker", lis.Addr().String())
	defer registryClient.Unregister("key_maker", lis.Addr().String())
	grpcServer.Serve(lis)
}

func (s *keyMakerServerImplementation) CreateUserFromKey(ctx context.Context, in *service.CreateUserKeyRequest) (*empty.Empty, error) {
	err := CreateUser(in.UserId, in.UserEncryptionKey)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *keyMakerServerImplementation) GetUserEncryptedKey(ctx context.Context, in *service.GetUserEncryptedKeyRequest) (*wrapperspb.StringValue, error) {
	encryptedKey, err := GetUserEncryptedKey(in.UserId)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{Value: encryptedKey}, nil
}

func (s *keyMakerServerImplementation) ChangeUserKey(ctx context.Context, in *service.ChangeUserKeyRequest) (*emptypb.Empty, error) {
	err := ChangeUserKey(in.UserId, in.OldKey, in.NewKey)
	return &emptypb.Empty{}, err
}
