package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"sgublogsite/src/internal/controller"
)

func registerHandlers() http.Handler {
    e := echo.New()
    useMiddleware(e)
    registerRoutes(e)
    return e
}

func useMiddleware(e *echo.Echo) {
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins:     []string{"https://*", "http://*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
        AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        AllowCredentials: true,
        MaxAge:           300,
    }))
}

func registerRoutes(e *echo.Echo) http.Handler {

    // Authorization
    e.POST("/register"   , controllers.Register)
    e.POST("/login"      , controllers.Login)
    e.POST("/logout"     , controllers.Logout)
    e.POST("/protected"  , controllers.Protected)

    return e
}
