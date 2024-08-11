package Server

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDatabaseConnection() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		return nil, err
	}
	return dbpool, nil
}

func GetUserEncryptedKey(userId string) (string, error) {
	dbpool, err := CreateDatabaseConnection()
	if err != nil {
		return "", err
	}
	defer dbpool.Close()
	var encryptedKey []byte
	err = dbpool.QueryRow(context.Background(), "SELECT userKEK FROM users WHERE userid = $1", userId).Scan(&encryptedKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return "", err
	}
	encryptedKeyHex := hex.EncodeToString(encryptedKey)
	fmt.Printf("Encrypted key: %s\n", encryptedKeyHex)
	return encryptedKeyHex, nil
}

func CreateUser(userId string, userEncryptionKey string) error {
	dbpool, err := CreateDatabaseConnection()
	if err != nil {
		return err
	}
	defer dbpool.Close()

	userKEK, err := GenerateNewKey()
	if err != nil {
		fmt.Printf("Error generating new key: %v\n", err)
		return err
	}

	encryptedKEK, err := EncryptKeyFromUserKey(userEncryptionKey, userKEK)
	if err != nil {
		fmt.Printf("Error encrypting key: %v\n", err)
		return err
	}

	_, err = dbpool.Exec(context.Background(), "UPDATE users SET userKEK = $1 WHERE userid = $2", encryptedKEK, userId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return err
	}

	fmt.Printf("User created: %s\n", userId)
	return nil
}

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

func EncryptKeyFromUserKey(userKey string, KEK string) ([]byte, error) {
	userKeyBytes, err := hex.DecodeString(userKey)
	if err != nil {
		return nil, err
	}
	KEKbytes, err := hex.DecodeString(KEK)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(userKeyBytes)
	if err != nil {
		log.Fatalf("cipher err: %v", err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("cipher GCM err: %v", err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("nonce  err: %v", err.Error())
	}

	fmt.Printf("Nonce: %v\n", hex.EncodeToString(nonce))
	encyptedKEK := gcm.Seal(nonce, nonce, KEKbytes, nil)
	return encyptedKEK, nil
}

func DecryptKeyFromUserKey(userKey []byte, KEK []byte) ([]byte, error) {
	block, err := aes.NewCipher(userKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Nonce Size: %v\n", gcm.NonceSize())

	nonce := KEK[:gcm.NonceSize()]
	cipherText := KEK[gcm.NonceSize():]
	decryptedKEK, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return decryptedKEK, nil
}

func ChangeUserKey(userId string, oldKeyHex string, newKeyHex string) error {
	oldKeyBytes, err := hex.DecodeString(oldKeyHex)
	if err != nil {
		return err
	}

	dbpool, err := CreateDatabaseConnection()
	if err != nil {
		return err
	}

	defer dbpool.Close()

	var encryptedKey []byte
	err = dbpool.QueryRow(context.Background(), "SELECT userKEK FROM users WHERE userid = $1", userId).Scan(&encryptedKey)
	if err != nil {
		return err
	}

	decryptedKey, err := DecryptKeyFromUserKey(oldKeyBytes, encryptedKey)
	if err != nil {
		return err
	}

	decryptedKeyHex := hex.EncodeToString(decryptedKey)

	encryptedNewKey, err := EncryptKeyFromUserKey(newKeyHex, decryptedKeyHex)
	if err != nil {
		return err
	}

	_, err = dbpool.Exec(context.Background(), "UPDATE users SET userKEK = $1 WHERE userid = $2", encryptedNewKey, userId)
	if err != nil {
		return err
	}

	return nil
}
