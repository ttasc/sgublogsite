package controller

import (
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

    posts, _         := c.Model.GetPostsByCategorySlug(
        "news",
        postsLimitPerPage,
        0,
        string(repos.PostsStatusPublished),
        isAuthenticated,
    )
    announcements, _ := c.Model.GetPostsByCategorySlug(
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

    if r.Header.Get("HX-Request") == "true" {
        c.templates["home"].ExecuteTemplate(w, "content", data)
    } else {
        c.templates["home"].Execute(w, data)
    }
}
