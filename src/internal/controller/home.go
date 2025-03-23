package controller

import (
	"html/template"
	"net/http"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"

	"github.com/go-chi/jwtauth/v5"
)

func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {
    isAuthenticated := false
    _, claims, err := jwtauth.FromContext(r.Context())
    if claims != nil && err == nil {
        isAuthenticated = true
    }

    posts, _         := c.model.GetPostsByCategorySlug(
        "news",
        postsLimitPerPage,
        0,
        string(repos.PostsStatusPublished),
        isAuthenticated,
    )
    announcements, _ := c.model.GetPostsByCategorySlug(
        "announcements",
        postsLimitPerPage,
        0,
        string(repos.PostsStatusPublished),
        isAuthenticated,
    )
    data := struct {
        IsAuthenticated bool
        Posts         []repos.GetPostsByCategorySlugRow
        Announcements []repos.GetPostsByCategorySlugRow
    }{
        IsAuthenticated: isAuthenticated,
        Posts:           posts,
        Announcements:   announcements,
    }
    tmpl, err := template.Must(c.basetmpl.Clone()).ParseFiles("templates/home.tmpl")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if r.Header.Get("HX-Request") == "true" {
        tmpl.ExecuteTemplate(w, "content", data)
    } else {
        tmpl.Execute(w, data)
    }
}
