package Client

import (
	"context"
	"fmt"
	"math/rand"

	RegistryClient "github.com/DeebTibi/GoVault/services/registry/client"
	service "github.com/DeebTibi/GoVault/services/user_auth/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserAuthClient struct {
	registryClient       *RegistryClient.RegistryClient
	createUserAuthClient func(cc grpc.ClientConnInterface) service.UserAuthClient
}

func NewUserAuthClient() *UserAuthClient {
	return &UserAuthClient{
		registryClient:       RegistryClient.NewRegistryClient(),
		createUserAuthClient: service.NewUserAuthClient,
	}
}

func (obj *UserAuthClient) PickRandomKeyMakerServer() (string, error) {
	ips, err := obj.registryClient.Discover("user_auth")
	fmt.Printf("Selected IP: %v\n", ips)
	if err != nil {
		return "", err
	}

	randomNumber := rand.Intn(len(ips))

	selectedIP := ips[randomNumber]
	return selectedIP, nil

}

func (obj *UserAuthClient) Login(userId string, password string) (string, error) {
	ip, err := obj.PickRandomKeyMakerServer()
	if err != nil {
		return "", err
	}

	conn, err := grpc.NewClient(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := obj.createUserAuthClient(conn)
	response, err := client.Login(context.Background(), &service.LoginRequest{UserName: userId, Password: password})
	if err != nil {
		return "", err
	}
	return response.Token, nil
}

func (obj *UserAuthClient) Register(userId string, password string, userKey string) (string, error) {
	ip, err := obj.PickRandomKeyMakerServer()
	if err != nil {
		return "", err
	}

	conn, err := grpc.NewClient(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := obj.createUserAuthClient(conn)
	response, err := client.Register(context.Background(), &service.RegisterRequest{UserName: userId, Password: password, UserKey: userKey})
	if err != nil {
		return "", err
	}
	return response.Token, nil
}

func (obj *UserAuthClient) Authenticate(userId string, token string) (bool, error) {
	ip, err := obj.PickRandomKeyMakerServer()
	if err != nil {
		return false, err
	}

	conn, err := grpc.NewClient(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := obj.createUserAuthClient(conn)
	response, err := client.AuthenticateToken(context.Background(), &service.AuthenticateTokenRequest{UserName: userId, Token: token})
	if err != nil {
		return false, err
	}
	return response.Value, nil
}
