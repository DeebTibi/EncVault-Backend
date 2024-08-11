package Server

import (
	"crypto/rand"
	"encoding/hex"
)

var userTokens = make(map[string]string)

func CreateUserToken(userId string) (string, error) {

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	userToken := hex.EncodeToString(b)
	userTokens[userId] = userToken
	return userToken, nil
}

func ValidateUserToken(userId string, token string) bool {
	// validate the token
	if userToken, ok := userTokens[userId]; ok {
		return userToken == token
	}
	return false
}
