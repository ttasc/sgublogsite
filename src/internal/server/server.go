package server

import (
    "fmt"
    "net/http"
    "os"
    "strconv"
    "time"

    _ "github.com/joho/godotenv/autoload"
)

func NewServer() *http.Server {
    port, _ := strconv.Atoi(os.Getenv("PORT"))

    // Declare Server config
    server := &http.Server{
        Addr:         fmt.Sprintf(":%d", port),
        Handler:      registerHandlers(),
        IdleTimeout:  time.Minute,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
    }

    return server
}
