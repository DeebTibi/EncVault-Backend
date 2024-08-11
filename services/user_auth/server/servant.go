package Server

import (
	"context"
	"fmt"
	"log"
	"os"

	"crypto/sha256"

	KeyMakerClient "github.com/DeebTibi/GoVault/services/key_maker/client"
	TokenGeneratorClient "github.com/DeebTibi/GoVault/services/token_generator/client"
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

func ComparePasswords(slice1, slice2 []byte) bool {
	// Check if lengths are the same
	if len(slice1) != len(slice2) {
		return false
	}

	// Compare elements one by one
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}

func HashPassword(password string) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write([]byte(password))
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func LoginUser(username string, password string) (string, error) {
	dbpool, err := CreateDatabaseConnection()
	if err != nil {
		return "", err
	}
	defer dbpool.Close()

	inputPassword, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	row := dbpool.QueryRow(context.Background(), "SELECT password_hash FROM users WHERE userid = $1", username)

	var dbPassword []byte

	err = row.Scan(&dbPassword)
	if err != nil {
		log.Fatalf("QueryRow failed: %v", err)
	}

	if !ComparePasswords(dbPassword, inputPassword) {
		return "", fmt.Errorf("invalid input password")
	}

	// generate a token for the user
	client := TokenGeneratorClient.NewTokenGeneratorClient()
	token, err := client.CreateUserToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func RegisterUser(username string, password string, userKey string) (string, error) {
	dbpool, err := CreateDatabaseConnection()
	if err != nil {
		return "", err
	}
	defer dbpool.Close()

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	// make sure client doesnt already exist
	_, err = dbpool.Exec(context.Background(), `INSERT INTO users (userid, password_hash, userKEK) VALUES ($1, $2, $3)`, username, hashedPassword, []byte("example"))
	if err != nil {
		return "", err
	}

	// generate and encrypt key using userkey
	kmClient := KeyMakerClient.NewKeyMakerClient()
	err = kmClient.CreateUserFromKey(username, userKey)
	if err != nil {
		return "", err
	}

	// generate a token for the user
	client := TokenGeneratorClient.NewTokenGeneratorClient()
	token, err := client.CreateUserToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func AuthenticateUser(userId string, token string) (bool, error) {
	client := TokenGeneratorClient.NewTokenGeneratorClient()
	res, err := client.ValidateToken(userId, token)
	if err != nil {
		return false, err
	}
	return res, nil
}
