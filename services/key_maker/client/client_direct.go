package Client

import (
	"context"

	ClientDirect "github.com/DeebTibi/GoVault/services/common"
	service "github.com/DeebTibi/GoVault/services/key_maker/api"
)

type KeyMakerClientDirect struct {
	ClientDirect.ClientDirect[service.KeyMakerClient]
}

func NewKeyMakerClientDirect(ip string) *KeyMakerClientDirect {
	return &KeyMakerClientDirect{
		ClientDirect.ClientDirect[service.KeyMakerClient]{
			Address:      ip,
			CreateClient: service.NewKeyMakerClient,
		},
	}
}

func (obj *KeyMakerClientDirect) CreateUserFromKey(userId string, userEncrpytionKey string) error {
	disconnect, client, err := obj.Connect()
	if err != nil {
		return err
	}
	defer disconnect()

	_, err = client.CreateUserFromKey(context.Background(), &service.CreateUserKeyRequest{UserId: userId, UserEncryptionKey: userEncrpytionKey})
	if err != nil {
		return err
	}
	return nil
}
