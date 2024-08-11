package Server

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/DeebTibi/GoVault/config"
	utils "github.com/DeebTibi/GoVault/services/file_upload/server/utils"
	RegistryClient "github.com/DeebTibi/GoVault/services/registry/client"
	"github.com/gofiber/fiber/v2"
)

func Start(cfg *config.ServiceConfig) {
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", 0755)
	}

	app := fiber.New()

	// Attach the authentication middleware to the /upload route
	app.Post("/upload", utils.Authenticate, func(c *fiber.Ctx) error {
		// Get the User-ID from the headers
		userID := c.Get("User-ID")
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing User-ID header")
		}

		userKey := c.Get("User-Key")
		if userKey == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing User-Key header")
		}

		// Create the user's directory if it doesn't exist
		userDir := fmt.Sprintf("./uploads/%s", userID)
		if _, err := os.Stat(userDir); os.IsNotExist(err) {
			os.Mkdir(userDir, 0755)
		}

		// Get the uploaded file from the form
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving file")
		}

		// Open the uploaded file
		fileContent, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error opening file")
		}
		defer fileContent.Close()

		// Read the file content into a byte slice
		fileBytes, err := io.ReadAll(fileContent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading file")
		}

		encryptedFile, err := utils.EncryptFile(userID, userKey, fileBytes)
		if err != nil {
			fmt.Printf("Error encrypting file: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error encrypting file")
		}

		// Save the encrypted file to the user's directory
		err = os.WriteFile(fmt.Sprintf("%s/%s", userDir, file.Filename), encryptedFile, 0644)
		if err != nil {
			fmt.Printf("Error saving file: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving file")
		}

		return c.SendString(fmt.Sprintf("File uploaded successfully: %s\n", file.Filename))
	})

	app.Get("/myfiles", utils.Authenticate, func(c *fiber.Ctx) error {
		// Get the User-ID from the headers
		userID := c.Get("User-ID")
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing User-ID header")
		}

		// Get the user's directory
		userDir := fmt.Sprintf("./uploads/%s", userID)
		// Get all file names from the directory
		files, err := os.ReadDir(userDir)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading files")
		}
		// return the list of files
		fileList := ""
		for _, file := range files {
			fileList += file.Name() + "\n"
		}
		return c.SendString(fileList)
	})

	app.Get("/download/", utils.Authenticate, func(c *fiber.Ctx) error {
		// Get the User-ID from the headers
		userID := c.Get("User-ID")
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing User-ID header")
		}

		userKey := c.Get("User-Key")
		if userKey == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing User-Key header")
		}

		// Get the file name from the URL parameter
		fileName := c.Query("file")
		if fileName == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing file name")
		}

		// Get the user's directory
		userDir := fmt.Sprintf("./uploads/%s", userID)
		// Read the encrypted file
		encryptedFile, err := os.ReadFile(fmt.Sprintf("%s/%s", userDir, fileName))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading file")
		}

		decryptedFile, err := utils.DecryptFile(userID, userKey, encryptedFile)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error decrypting file")
		}

		// Determine the MIME type
		mimeType := mime.TypeByExtension(filepath.Ext(fileName))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
		// Set the Content-Type header
		c.Set(fiber.HeaderContentType, mimeType)

		reader := bytes.NewReader(decryptedFile)

		// Return the decrypted file as a download
		return c.SendStream(io.NopCloser(reader))
	})

	regClient := RegistryClient.NewRegistryClient()
	regClient.Register("file_upload", "localhost:8080")
	defer regClient.Unregister("file_upload", "localhost:8080")

	fmt.Println("Server started at :8080")
	app.Listen(":8080")
}
