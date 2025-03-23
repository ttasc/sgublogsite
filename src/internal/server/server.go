package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/ttasc/sgublogsite/src/internal/controller"
	"github.com/ttasc/sgublogsite/src/internal/model"
	"github.com/ttasc/sgublogsite/src/internal/utils"
)

type Server struct {
    port    int
    ctrlr   controller.Controller
}

func NewServer() *http.Server {
    port, _ := strconv.Atoi(os.Getenv("PORT"))
    NewServer := &Server{
        port:    port,
        ctrlr:   controller.New(model.New(utils.NewDB())),
    }

    // Declare Server config
    server := &http.Server{
        Addr:         fmt.Sprintf(":%d", NewServer.port),
        Handler:      NewServer.registerHandlers(),
        IdleTimeout:  time.Minute,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
    }

    return server
}
