package server

import (
	"sgublogsite/src/internal/controller"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func userMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        user := controller.User{IsAuthenticated: false}
        if token, ok := c.Get("jwt").(*jwt.Token); ok && token.Valid {
            if claims, ok := token.Claims.(*controller.JWTClaims); ok {
                user = controller.User{
                    ID: claims.ID,
                    Roles: claims.Roles,
                    IsAuthenticated: true,
                }
            }
        }
        c.Set("user", user)
        return next(c)
    }
}

func configureJWT() echo.MiddlewareFunc {
    return echojwt.WithConfig(echojwt.Config{
        SigningKey:  controller.JWTKey,
        TokenLookup: "cookie:" + controller.JWTCookieName,
        ContextKey:  "jwt",
        NewClaimsFunc: func(c echo.Context) jwt.Claims {
            return new(controller.JWTClaims)
        },
        ErrorHandler: func(c echo.Context, err error) error {
            return nil
        },
    })
}
