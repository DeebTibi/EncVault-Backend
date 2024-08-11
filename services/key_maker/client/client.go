package Client

import (
	"context"
	"math/rand"

	service "github.com/DeebTibi/GoVault/services/key_maker/api"
	RegistryClient "github.com/DeebTibi/GoVault/services/registry/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type KeyMakerClient struct {
	registryClient       *RegistryClient.RegistryClient
	createKeyMakerClient func(cc grpc.ClientConnInterface) service.KeyMakerClient
}

func NewKeyMakerClient() *KeyMakerClient {
	return &KeyMakerClient{
		registryClient:       RegistryClient.NewRegistryClient(),
		createKeyMakerClient: service.NewKeyMakerClient,
	}
}

func (obj *KeyMakerClient) PickRandomKeyMakerServer() (string, error) {
	ips, err := obj.registryClient.Discover("key_maker")
	if err != nil {
		return "", err
	}

	randomNumber := rand.Intn(len(ips))

	selectedIP := ips[randomNumber]
	return selectedIP, nil

}

func (obj *KeyMakerClient) CreateUserFromKey(userId string, userEncrpytionKey string) error {
	ip, err := obj.PickRandomKeyMakerServer()
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := obj.createKeyMakerClient(conn)
	_, err = client.CreateUserFromKey(context.Background(), &service.CreateUserKeyRequest{UserId: userId, UserEncryptionKey: userEncrpytionKey})
	if err != nil {
		return err
	}
	return nil
}

func (obj *KeyMakerClient) GetUserEncryptedKey(userId string) (string, error) {
	ip, err := obj.PickRandomKeyMakerServer()
	if err != nil {
		return "", err
	}

	conn, err := grpc.NewClient(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := obj.createKeyMakerClient(conn)
	response, err := client.GetUserEncryptedKey(context.Background(), &service.GetUserEncryptedKeyRequest{UserId: userId})
	if err != nil {
		return "", err
	}
	return response.Value, nil
}

func (obj *KeyMakerClient) ChangeUserKey(userId string, oldKey string, newKey string) error {
	ip, err := obj.PickRandomKeyMakerServer()
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := obj.createKeyMakerClient(conn)
	_, err = client.ChangeUserKey(context.Background(), &service.ChangeUserKeyRequest{UserId: userId, OldKey: oldKey, NewKey: newKey})
	if err != nil {
		return err
	}
	return nil
}
