package Client

import (
	"context"
	"math/rand"

	RegistryClient "github.com/DeebTibi/GoVault/services/registry/client"
	service "github.com/DeebTibi/GoVault/services/token_generator/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TokenGeneratorClient struct {
	registryClient             *RegistryClient.RegistryClient
	createTokenGeneratorClient func(cc grpc.ClientConnInterface) service.TokenGeneratorClient
}

func NewTokenGeneratorClient() *TokenGeneratorClient {
	return &TokenGeneratorClient{
		registryClient:             RegistryClient.NewRegistryClient(),
		createTokenGeneratorClient: service.NewTokenGeneratorClient,
	}
}

func (obj *TokenGeneratorClient) PickRandomTokenGeneratorServer() (string, error) {
	ips, err := obj.registryClient.Discover("token_generator")
	if err != nil {
		return "", err
	}

	randomNumber := rand.Intn(len(ips))

	selectedIP := ips[randomNumber]
	return selectedIP, nil
}

func (obj *TokenGeneratorClient) CreateUserToken(userId string) (string, error) {
	ip, err := obj.PickRandomTokenGeneratorServer()
	if err != nil {
		return "", err
	}

	conn, err := grpc.NewClient(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := obj.createTokenGeneratorClient(conn)
	token, err := client.CreateUserToken(context.Background(), &service.CreateUserTokenRequest{UserId: userId})
	if err != nil {
		return "", err
	}
	return token.Value, nil
}

func (obj *TokenGeneratorClient) ValidateToken(userId string, token string) (bool, error) {
	ip, err := obj.PickRandomTokenGeneratorServer()
	if err != nil {
		return false, err
	}

	conn, err := grpc.NewClient(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := obj.createTokenGeneratorClient(conn)
	res, err := client.ValidateUserToken(context.Background(), &service.ValidateUserTokenRequest{UserId: userId, UserToken: token})
	if err != nil {
		return false, err
	}
	return res.Value, nil
}
