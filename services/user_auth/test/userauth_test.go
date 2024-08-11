package Tester

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	clientUserAuth "github.com/DeebTibi/GoVault/services/user_auth/client"
)

func GenerateNewKey() (string, error) {
	// generate a 16 byte key
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	keyHex := hex.EncodeToString(key)
	return keyHex, nil
}

func TestCreateUser(t *testing.T) {
	c := clientUserAuth.NewUserAuthClient()
	key, err := GenerateNewKey()
	fmt.Printf("UserKey: %s\n", key)
	if err != nil {
		t.Error(err)
	}
	token, err := c.Register("deeb5", "password", key)
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("Token is empty")
	}
	fmt.Printf("Token: %s\n", token)
}
