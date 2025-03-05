package server

import (
	"net/http"
	controller "website/src/internal/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func registerRoutes() http.Handler {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins:     []string{"https://*", "http://*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
        AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    e.POST("/register"   , controller.Register)
    e.POST("/login"      , controller.Login)
    e.POST("/logout"     , controller.Logout)
    e.POST("/protected"  , controller.Protected)

    return e
}

