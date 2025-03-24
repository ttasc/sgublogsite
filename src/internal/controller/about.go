package controller

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func (c *Controller) About(w http.ResponseWriter, r *http.Request) {
    about, _ := c.Model.GetSiteAbout()
    data := struct {
        IsAuthenticated bool
        About string
    }{
        About: about,
    }
    if r.Header.Get("HX-Request") == "true" {
        c.templates["about"].ExecuteTemplate(w, "content", data)
    } else {
        _, claims, err := jwtauth.FromContext(r.Context())
        data.IsAuthenticated = (claims != nil && err == nil)
        c.templates["about"].Execute(w, data)
    }
}
