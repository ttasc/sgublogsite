package controller

import (
	"html/template"
	"net/http"
	"sgublogsite/src/internal/model"
	"sgublogsite/src/internal/model/repos"

	"github.com/go-chi/jwtauth/v5"
)

func Home(w http.ResponseWriter, r *http.Request) {
    isAuthenticated := false
    _, claims, err := jwtauth.FromContext(r.Context())
    if claims != nil && err == nil {
        isAuthenticated = true
    }

    m := model.New()
    posts, _         := m.GetPostsByCategorySlug(
        "news",
        postsLimitPerPage,
        0,
        string(repos.PostsStatusPublished),
        isAuthenticated,
    )
    announcements, _ := m.GetPostsByCategorySlug(
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
    tmpl, err := template.Must(basetmpl.Clone()).ParseFiles("templates/home.tmpl")
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
