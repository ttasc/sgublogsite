package controller

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
	_ "github.com/joho/godotenv/autoload"

	"github.com/ttasc/sgublogsite/src/internal/model"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
	// "github.com/ttasc/sgublogsite/src/internal/utils"
)

var (
    JWTCookieName   = "jwt"
    expirationTime  = time.Now().Add(time.Hour * 24)
    jwtKey          = []byte(os.Getenv("SGUBLOGSITE_JWT_KEY"))
    TokenAuth       *jwtauth.JWTAuth
)

func init() {
    TokenAuth = jwtauth.New("HS256", jwtKey, nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("statics/login.html"))
    tmpl.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte("Method not allowed"))
        return
    }

    emailorphone := r.FormValue("emailorphone")
    // password := r.FormValue("password")

    m := model.New()
    user, err := m.GetUserByEmailOrPhone(emailorphone)

    // if err != nil || !utils.CheckPasswordHash(password, user.Password) {
    //     w.WriteHeader(http.StatusUnauthorized)
    //     w.Write([]byte("Invalid email or password"))
    //     return
    // }

    _, tokenString, err := TokenAuth.Encode(map[string]any{
        "ID":    user.UserID,
        "Roles": []string{string(user.Role.UsersRole)},
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

    if user.Role.UsersRole == repos.UsersRoleAdmin {
        http.Redirect(w, r, "/admin", http.StatusSeeOther)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
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
