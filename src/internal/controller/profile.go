package controller

import (
	"html/template"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func (c *Controller) Profile(w http.ResponseWriter, r *http.Request) {
    _, claims, err := jwtauth.FromContext(r.Context())
    if claims == nil || err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    currentUser, _ := c.Model.GetUserByID(int32(claims["ID"].(float64)))

    tmpl := template.Must(template.ParseFiles("templates/profile.tmpl"))
    tmpl.Execute(w, currentUser)
}
