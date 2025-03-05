package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"website/internal/database"
	"website/internal/utils"
)

func Protected(c echo.Context) error {
    if c.Request().Method != http.MethodPost {
        c.String(http.StatusMethodNotAllowed, "Invalid request method")
        return errors.New("Invalid request method")
    }

    if err := authorize(c); err != nil {
        c.String(http.StatusUnauthorized, "Unauthorized")
        return errors.New("Unauthorized")
    }

    email := c.FormValue("email")
    c.String(http.StatusOK, "Protected route for "+email)

    return nil
}

func Register(c echo.Context) error {
    if c.Request().Method != http.MethodPost {
        c.String(http.StatusMethodNotAllowed, "Method not allowed")
        return errors.New("Method not allowed")
    }
    email := c.FormValue("email")
    fullname := c.FormValue("fullname")
    password := c.FormValue("password")
    if len(email)==0 || len(password)<8 {
        c.String(http.StatusBadRequest, "Invalid email or password")
        return errors.New("Invalid email or password")
    }
    if database.UserExists(email) {
        c.String(http.StatusConflict, "User already exists")
        return errors.New("User already exists")
    }
    hashedPassword, _ := utils.HashPassword(password)
    user := &database.User{
        FullName: fullname,
        Password: hashedPassword,
        Email:    email,
    }
    database.AddUser(user)

    c.String(http.StatusCreated, "User created successfully")

    return nil
}

func Login(c echo.Context) error {
    if c.Request().Method != http.MethodPost {
        c.String(http.StatusMethodNotAllowed, "Method not allowed")
        return errors.New("Method not allowed")
    }
    email := c.FormValue("email")
    password := c.FormValue("password")
    user, err := database.GetUserByEmail(email)
    if err != nil || !utils.CheckPasswordHash(password, user.Password) {
        c.String(http.StatusBadRequest, "Invalid email or password")
        return errors.New("Invalid email or password")
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
        c.String(http.StatusInternalServerError, "Failed to add session data")
        return errors.New("Failed to add session data")
    }

    c.String(http.StatusOK, "Login successfully")

    return nil
}

func Logout(c echo.Context) error {
    if c.Request().Method != http.MethodPost {
        c.String(http.StatusMethodNotAllowed, "Invalid request method")
        return errors.New("Invalid request method")
    }
    if err := authorize(c); err != nil {
        c.String(http.StatusUnauthorized, "Unauthorized")
        return errors.New("Unauthorized")
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

    // Clear the tokens from the database
    email := c.FormValue("email")
    user, err := database.GetUserByEmail(email)
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
    user, err := database.GetUserByEmail(c.FormValue("email"))
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
