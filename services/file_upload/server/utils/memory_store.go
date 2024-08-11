package utils

import (
	"encoding/hex"
	"fmt"
)

var userKeys = make(map[string]string)

func GetUserKey(userID string) ([]byte, error) {
	key, ok := userKeys[userID]
	if !ok {
		return nil, fmt.Errorf("key not found for user %s", userID)
	}

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key for user %s: %v", userID, err)
	}

	return keyBytes, nil
}

func SetUserKey(userID string, key []byte) {
	keyString := hex.EncodeToString(key)
	userKeys[userID] = keyString
}
