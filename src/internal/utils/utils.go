package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func GenerateToken(length int) string {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        log.Fatalf("Failed to generate token: %v", err)
        return ""
    }
    return base64.URLEncoding.EncodeToString(bytes)
}

func SaveUploadedFile(file multipart.File, handler *multipart.FileHeader, prefix_path string) (string, error) {
    defer file.Close()

    var uploadABDir = prefix_path
    var uploadREDIR = fmt.Sprintf(".%s", uploadABDir)

    // Create upload directory if it doesn't exist
    if err := os.MkdirAll(uploadREDIR, 0755); err != nil {
        return "", err
    }

    filename := handler.Filename
    filePath := path.Join(uploadREDIR, filename)

    // Write the file to disk
    f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        return "", err
    }
    defer f.Close()

    io.Copy(f, file)
    return path.Join(uploadABDir, filename), nil
}

func DeleteUploadedFile(url string) error {
    if url == "" {
        return nil
    }

    filePath := fmt.Sprintf(".%s", url)

    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return fmt.Errorf("file not found: %s", filePath)
    }

    if err := os.Remove(filePath); err != nil {
        return fmt.Errorf("failed to delete file: %v", err)
    }

    return nil
}

func GenerateUniqueFilename(basePath, originalName string) string {
    ext := path.Ext(originalName)
    name := strings.TrimSuffix(originalName, ext)
    counter := 0

    for {
        newName := originalName
        if counter > 0 {
            newName = fmt.Sprintf("%s(%d)%s", name, counter, ext)
        }

        fullPath := path.Join(basePath, newName)
        if _, err := os.Stat(fullPath); os.IsNotExist(err) {
            return newName
        }
        counter++
    }
}
