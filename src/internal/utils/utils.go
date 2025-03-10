package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"

	"github.com/labstack/echo/v4"
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

func Rdie(c echo.Context, statusCode int, message string) error {
    c.String(statusCode, message)
    return errors.New(message)
}
