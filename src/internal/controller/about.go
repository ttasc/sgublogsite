package controller

import (
	"html/template"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func (c *Controller) About(w http.ResponseWriter, r *http.Request) {
    about, _ := c.model.GetSiteAbout()
    data := struct {
        IsAuthenticated bool
        About string
    }{
        About: about,
    }
    tmpl, err := template.Must(c.basetmpl.Clone()).ParseFiles("templates/about.tmpl")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if r.Header.Get("HX-Request") == "true" {
        tmpl.ExecuteTemplate(w, "content", data)
    } else {
        _, claims, err := jwtauth.FromContext(r.Context())
        data.IsAuthenticated = (claims != nil && err == nil)
        tmpl.Execute(w, data)
    }
}
