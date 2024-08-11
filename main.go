package main

import (
	"fmt"
	"os"

	"github.com/DeebTibi/GoVault/config"
	FileUploadServer "github.com/DeebTibi/GoVault/services/file_upload/server"
	KeyMakerServer "github.com/DeebTibi/GoVault/services/key_maker/server"
	RegistryServer "github.com/DeebTibi/GoVault/services/registry/server"
	TokenGeneratorServer "github.com/DeebTibi/GoVault/services/token_generator/server"
	UserAuthServer "github.com/DeebTibi/GoVault/services/user_auth/server"
	"gopkg.in/yaml.v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <config file>")
		return
	}

	configFile := os.Args[1]
	fileData, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	cfg := &config.ServiceConfig{}
	err = yaml.Unmarshal(fileData, cfg)
	if err != nil {
		fmt.Printf("Error unmarshaling YAML: %v\n", err)
		return
	}

	fmt.Printf("ServiceName: %s\n", cfg.ServiceName)

	switch cfg.ServiceName {
	case "key_maker":
		fmt.Println("Starting key_maker service")
		KeyMakerServer.Start(cfg)
	case "file_upload":
		fmt.Println("Starting file_upload service")
		FileUploadServer.Start(cfg)
	case "registry":
		fmt.Println("Starting registry service")
		RegistryServer.Start(cfg)
	case "token_generator":
		fmt.Println("Starting token_generator service")
		TokenGeneratorServer.Start(cfg)
	case "user_auth":
		fmt.Println("Starting user_auth service")
		UserAuthServer.Start(cfg)
	default:
		fmt.Println("Unknown service")
	}

}
