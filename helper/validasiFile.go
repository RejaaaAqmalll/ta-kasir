package helper

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

func IsSupportedImageFormat(mime string) bool {
	return strings.HasPrefix(mime, "image/jpeg") || strings.HasPrefix(mime, "image/png") || strings.HasPrefix(mime, "image/jpg")
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	filename := fmt.Sprintf("%s%s", GenerateRandomString(10), ext)
	return filename
}

func SaveFile(src io.Reader, filename string) error {
	// Specify your storage directory
	storageDir := "./storage/foto"

	// Create the storage directory if it doesn't exist
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		return err
	}

	 // Check if the file already exists
	 filePath := filepath.Join(storageDir, filename)
	 if _, err := os.Stat(filePath); err == nil {
		 return fmt.Errorf("file already exists: %s", filename)
	 }

	// Create the file on the storage directory
    dst, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("failed to create destination file: %v", err)
    }
    defer dst.Close()

	 // Copy the contents from src to dst
	 _, err = io.Copy(dst, src)
	 if err != nil {
		 return fmt.Errorf("failed to copy contents to destination file: %v", err)
	 }

	log.Printf("File saved successfully: %s\n", filename)
	return nil
}