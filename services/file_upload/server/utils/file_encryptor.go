package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	KeyMakerClient "github.com/DeebTibi/GoVault/services/key_maker/client"
)

func EncryptFile(userId string, userKeyHex string, file []byte) ([]byte, error) {
	fmt.Printf("received key %v\n", userKeyHex)
	client := KeyMakerClient.NewKeyMakerClient()
	key, err := client.GetUserEncryptedKey(userId)
	if err != nil {
		return nil, err
	}
	fmt.Printf("User Encrypted key: %v\n", key)
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	userKey, err := hex.DecodeString(userKeyHex)
	if err != nil {
		return nil, err
	}

	decryptedKey, err := DecryptKEK(userKey, keyBytes)
	if err != nil {
		fmt.Printf("Error decrypting userkey\n")
		return nil, err
	}

	fmt.Printf("Decrypted user key: %v", hex.EncodeToString(decryptedKey))

	block, err := aes.NewCipher(decryptedKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("cipher GCM err: %v", err.Error())
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("nonce  err: %v", err.Error())
	}

	encryptedFile := gcm.Seal(nonce, nonce, file, nil)

	return encryptedFile, nil
}

func DecryptFile(userId string, userKeyHex string, file []byte) ([]byte, error) {

	client := KeyMakerClient.NewKeyMakerClient()
	keyHex, err := client.GetUserEncryptedKey(userId)
	if err != nil {
		return nil, err
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, err
	}

	userKey, err := hex.DecodeString(userKeyHex)
	if err != nil {
		return nil, err
	}

	decryptedKey, err := DecryptKEK(userKey, key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(decryptedKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("cipher GCM err: %v", err.Error())
	}

	nonce := file[:gcm.NonceSize()]
	cipherText := file[gcm.NonceSize():]
	decryptedFile, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Fatalf("decrypt file err: %v", err.Error())
	}

	return decryptedFile, nil
}

func DecryptKEK(userKey []byte, kek []byte) ([]byte, error) {
	block, err := aes.NewCipher(userKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Nonce Size: %v\n", gcm.NonceSize())

	nonce := kek[:gcm.NonceSize()]
	cipherText := kek[gcm.NonceSize():]
	decryptedKEK, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return decryptedKEK, nil
}
