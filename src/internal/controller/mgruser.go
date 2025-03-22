package controller

import (
	"net/http"
	"sgublogsite/src/internal/model"
	"sgublogsite/src/internal/model/repos"
	"sgublogsite/src/internal/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte("Method not allowed"))
        return
    }

    var (
        firstname   = r.FormValue("firstname")
        lastname    = r.FormValue("lastname")
        phone       = r.FormValue("phone")
        email       = r.FormValue("email")
        role        = r.FormValue("role")
        password    = r.FormValue("password")
    )
    m := model.New()

    hashedPassword, _ := utils.HashPassword(password)

    m.AddUser(repos.User{
        Firstname:   firstname,
        Lastname:    lastname,
        Phone:       phone,
        Email:       email,
        Password:    hashedPassword,
        Role:        repos.NullUsersRole{ UsersRole: repos.UsersRole(role), Valid: true },
    })

    w.Write([]byte("User registered successfully"))
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

