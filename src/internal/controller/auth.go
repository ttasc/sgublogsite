package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	_ "github.com/joho/godotenv/autoload"

	"sgublogsite/src/internal/model/repos"
	"sgublogsite/src/internal/model"
	"sgublogsite/src/internal/utils"
)

type JWTClaims struct {
    UserID string
    Roles []string `json:"role"`
    jwt.RegisteredClaims
}

type User struct {
    IsAuthenticated bool
    ID              string
    Roles           []string
}

var (
    JWTKey = []byte(os.Getenv("SGUBLOGSITE_JWT_KEY"))
    JWTCookieName  = "jwt_token"
)

func Register(c echo.Context) error {
    if c.Request().Method != http.MethodPost {
        return utils.Rdie(c, http.StatusMethodNotAllowed, "Method not allowed")
    }

    var (
        firstname   = c.FormValue("firstname")
        lastname    = c.FormValue("lastname")
        mobile      = c.FormValue("mobile")
        email       = c.FormValue("email")
        role        = c.FormValue("role")
        password    = c.FormValue("password")
    )
    m := model.New()

    hashedPassword, _ := utils.HashPassword(password)

    m.AddUser(repos.User{
        Firstname:   firstname,
        Lastname:    lastname,
        Mobile:      mobile,
        Email:       email,
        Password:    hashedPassword,
        Role:        repos.NullUsersRole{ UsersRole: repos.UsersRole(role), Valid: true },
    })

    return c.String(http.StatusCreated, "User created successfully")
}

func Login(c echo.Context) error {
    if c.Request().Method != http.MethodPost {
        return utils.Rdie(c, http.StatusMethodNotAllowed, "Method not allowed")
    }

    emailormobile := c.FormValue("emailormobile")
    password := c.FormValue("password")

    m := model.New()

    user, err := m.GetUserByEmailOrMobile(emailormobile)

    if err != nil || !utils.CheckPasswordHash(password, user.Password) {
        return utils.Rdie(c, http.StatusBadRequest, "Invalid email or password")
    }

    expirationTime := time.Now().Add(time.Hour * 24)

    // Set custom claims
    claims := &JWTClaims{
        UserID: string(user.UserID),
        Roles: []string{ string(user.Role.UsersRole) },
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    // Create token with claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Generate encoded token and send it as response.
    t, err := token.SignedString(JWTKey)
    if err != nil {
        return err
    }

    // Set cookie
    c.SetCookie(&http.Cookie{
        Name:     JWTCookieName,
        Value:    t,
        Expires:  expirationTime,
        HttpOnly: true,
        Secure:   false,
        Path:     "/",
    })

    return c.String(http.StatusOK, "Login successfully")
}

func Logout(c echo.Context) error {
    if c.Request().Method != http.MethodPost {
        return utils.Rdie(c, http.StatusMethodNotAllowed, "Method not allowed")
    }

    // Clear cookie
    c.SetCookie(&http.Cookie{
        Name:     "token",
        Value:    "",
        Expires:  time.Now().Add(-time.Hour),
        HttpOnly: true,
        Secure:   false,
        Path:     "/",
    })

    return c.String(http.StatusOK, "Logout successfully")
}
