package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ttasc/sgublogsite/src/internal/server"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
    // Create context that listens for the interrupt signal from the OS.
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    // Listen for the interrupt signal.
    <-ctx.Done()

    log.Println("shutting down gracefully, press Ctrl+C again to force")

    // The context is used to inform the server it has 5 seconds to finish
    // the request it is currently handling
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := apiServer.Shutdown(ctx); err != nil {
        log.Printf("Server forced to shutdown with error: %v", err)
    }

    log.Println("Server exiting")

    // Notify the main goroutine that the shutdown is complete
    done <- true
}

func main() {
    var err error
    // // Create a log file
    // file, err = os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    // if err != nil {
    //     log.Fatal("Failed to open log file:", err)
    // }
    // defer file.Close()
    // // Set output of logs to file
    // log.SetOutput(file)

    server := server.NewServer()

    // Create a done channel to signal when the shutdown is complete
    done := make(chan bool, 1)

    // Run graceful shutdown in a separate goroutine
    go gracefulShutdown(server, done)

    err = server.ListenAndServe()
    if err != nil && err != http.ErrServerClosed {
        panic(fmt.Sprintf("http server error: %s", err))
    }

    // Wait for the graceful shutdown to complete
    <-done
    log.Println("Graceful shutdown complete.")
}
