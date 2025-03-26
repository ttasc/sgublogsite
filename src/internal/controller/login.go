package controller

import (
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ttasc/sgublogsite/src/internal/model/repos"
	// "github.com/ttasc/sgublogsite/src/internal/utils"
)

var (
    JWTCookieName   = "jwt"
    expirationTime  = time.Now().Add(time.Hour * 24)
    jwtKey          = []byte(os.Getenv("SGUBLOGSITE_JWT_KEY"))
)

func (c *Controller) LoginPage(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "statics/login.html")
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte("Method not allowed"))
        return
    }

    emailorphone := r.FormValue("emailorphone")
    // password := r.FormValue("password")

    user, err := c.Model.GetUserByEmailOrPhone(emailorphone)

    // if err != nil || !utils.CheckPasswordHash(password, user.Password) {
    //     http.Redirect(w, r, "/login?error=auth_failed", http.StatusFound)
    //     return
    // }

    _, tokenString, err := c.TokenAuth.Encode(map[string]any{
        "ID":    user.UserID,
        "Roles": []string{string(user.Role)},
    })
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }
    // Set cookie
    http.SetCookie(w, &http.Cookie{
        Name:     JWTCookieName,
        Value:    tokenString,
        Expires:  expirationTime,
        HttpOnly: true,
        Secure:   false,
    })

    if user.Role == repos.UsersRoleAdmin {
        http.Redirect(w, r, "/admin", http.StatusSeeOther)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte("Method not allowed"))
        return
    }

    // Clear cookie
    http.SetCookie(w, &http.Cookie{
        Name:     JWTCookieName,
        Value:    "",
        Expires:  time.Now().Add(-time.Hour),
        HttpOnly: true,
        Secure:   false,
    })

    http.Redirect(w, r, "/", http.StatusSeeOther)
}
