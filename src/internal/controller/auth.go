package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"sgublogsite/src/internal/model/services"
	"sgublogsite/src/internal/utils"
)

func Protected(c echo.Context) error {
    if c.Request().Method != http.MethodPost {
        return utils.Rdie(c, http.StatusMethodNotAllowed, "Method not allowed")
    }

    if err := authorize(c); err != nil {
        return utils.Rdie(c, http.StatusUnauthorized, "Unauthorized")
    }

    email := c.FormValue("email")
    c.String(http.StatusOK, "Protected route for "+email)

    return nil
}

func Register(c echo.Context) error {
    if c.Request().Method != http.MethodPost {
        return utils.Rdie(c, http.StatusMethodNotAllowed, "Method not allowed")
    }
    email := c.FormValue("email")
    fullname := c.FormValue("fullname")
    password := c.FormValue("password")
    if len(email)==0 || len(password)<8 {
        return utils.Rdie(c, http.StatusBadRequest, "Invalid email or password")
    }
    if models.UserExists(email) {
        return utils.Rdie(c, http.StatusBadRequest, "User already exists")
    }
    hashedPassword, _ := utils.HashPassword(password)
    user := &models.User{
        FullName: fullname,
        Password: hashedPassword,
        Email:    email,
    }
    models.AddUser(user)

    c.String(http.StatusCreated, "User created successfully")

    return nil
}

func Login(c echo.Context) error {
    if c.Request().Method != http.MethodPost {
        return utils.Rdie(c, http.StatusMethodNotAllowed, "Method not allowed")
    }
    email := c.FormValue("email")
    password := c.FormValue("password")
    user, err := models.GetUserByEmail(email)
    if err != nil || !utils.CheckPasswordHash(password, user.Password) {
        return utils.Rdie(c, http.StatusBadRequest, "Invalid email or password")
    }

    sessionToken := utils.GenerateToken(32)
    csrfToken := utils.GenerateToken(32)

    // Set session cookie
    c.SetCookie(&http.Cookie{
        Name:     "session_token",
        Value:    sessionToken,
        Expires:  time.Now().Add(time.Hour * 24), // Set session expiration time
        HttpOnly: true,
    })

    // Set CSRF token cookie
    c.SetCookie(&http.Cookie{
        Name:     "csrf_token",
        Value:    csrfToken,
        Expires:  time.Now().Add(time.Hour * 24), // Set session expiration time
        HttpOnly: false, // Need to be accessible to the client-side
    })

    if err = user.UpdateSessionData(sessionToken, csrfToken); err != nil {
        return utils.Rdie(c, http.StatusInternalServerError, "Failed to update session data")
    }

    c.String(http.StatusOK, "Login successfully")

    return nil
}

func Logout(c echo.Context) error {
    if c.Request().Method != http.MethodPost {
        return utils.Rdie(c, http.StatusMethodNotAllowed, "Method not allowed")
    }
    if err := authorize(c); err != nil {
        return utils.Rdie(c, http.StatusUnauthorized, "Unauthorized")
    }

    // Clear cookie
    c.SetCookie(&http.Cookie{
        Name:     "session_token",
        Value:    "",
        Expires:  time.Now().Add(-time.Hour),
        HttpOnly: true,
    })

    c.SetCookie(&http.Cookie{
        Name:     "csrf_token",
        Value:    "",
        Expires:  time.Now().Add(-time.Hour),
        HttpOnly: false,
    })

    // Clear the tokens from the models
    email := c.FormValue("email")
    user, err := models.GetUserByEmail(email)
    if err != nil {
        return err
    }
    if err := user.UpdateSessionData("", ""); err != nil {
        return err
    }

    c.String(http.StatusOK, "Logout successfully")

    return nil
}

func authorize(c echo.Context) error {
    authError := errors.New("Unauthorized")
    user, err := models.GetUserByEmail(c.FormValue("email"))
    if err != nil {
        return authError
    }

    sessionTokenData, csrfTokenData, err := user.GetSessionData()
    if err != nil {
        return authError
    }

    // Get session token from cookie
    sessionToken, err := c.Cookie("session_token")
    if err != nil || sessionToken.Value == "" || sessionToken.Value != sessionTokenData {
        return authError
    }

    // Get CSRF token from cookie
    csrfToken := c.Request().Header.Get("X-CSRF-Token")
    if csrfToken == "" || csrfToken != csrfTokenData {
        return authError
    }

    return nil
}

