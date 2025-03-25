package utils

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"

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

func SaveUploadedFile(file multipart.File, handler *multipart.FileHeader) (string, error) {
    defer file.Close()

    // Create upload directory if it doesn't exist
    const uploadABDir = "./assets/uploads/avatars/"
    const uploadREDir = "/assets/uploads/avatars/"

    if err := os.MkdirAll(uploadABDir, 0755); err != nil {
        return "", err
    }

    filename := handler.Filename
    filePath := path.Join(uploadABDir, filename)

    // Write the file to disk
    f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        return "", err
    }
    defer f.Close()

    io.Copy(f, file)
    return path.Join(uploadREDir, filename), nil
}
