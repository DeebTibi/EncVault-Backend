package Client

import (
	"context"

	service "github.com/DeebTibi/GoVault/services/registry/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RegistryClient struct {
	registryIp   string
	CreateClient func(cc grpc.ClientConnInterface) service.RegistryClient
}

func NewRegistryClient() *RegistryClient {
	return &RegistryClient{
		registryIp:   "localhost:8502",
		CreateClient: service.NewRegistryClient,
	}
}

func (obj *RegistryClient) Discover(serviceName string) ([]string, error) {
	conn, err := grpc.NewClient(obj.registryIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := obj.CreateClient(conn)
	res, err := client.Discover(context.Background(), &service.DiscoverRequest{ServiceName: serviceName})
	if err != nil {
		return nil, err
	}
	return res.ServiceIps, nil
}

func (obj *RegistryClient) Register(serviceName, serviceIp string) error {
	conn, err := grpc.NewClient(obj.registryIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := obj.CreateClient(conn)
	_, err = client.Register(context.Background(), &service.RegisterRequest{ServiceName: serviceName, ServiceIp: serviceIp})
	if err != nil {
		return err
	}
	return nil
}

func (obj *RegistryClient) Unregister(serviceName, serviceIp string) error {
	conn, err := grpc.NewClient(obj.registryIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := obj.CreateClient(conn)
	_, err = client.Unregister(context.Background(), &service.UnregisterRequest{ServiceName: serviceName, ServiceIp: serviceIp})
	if err != nil {
		return err
	}
	return nil
}
