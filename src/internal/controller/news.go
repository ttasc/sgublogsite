package controller

import (
	"net/http"
	"strconv"

	"github.com/ttasc/sgublogsite/src/internal/model/repos"

	"github.com/go-chi/jwtauth/v5"
)

func (c *Controller) News(w http.ResponseWriter, r *http.Request) {
    isAuthenticated := false
    _, claims, err := jwtauth.FromContext(r.Context())
    if claims != nil && err == nil {
        isAuthenticated = true
    }

    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }
    offset := (int32) (page - 1) * postsLimitPerPage

    posts, _ := c.Model.GetPostsByCategorySlug(
        "news",
        postsLimitPerPage,
        offset,
        string(repos.PostsStatusPublished),
        isAuthenticated,
    )
    data := struct {
        IsAuthenticated bool
        Posts       []repos.GetPostsByCategorySlugRow
        Pagination  []paginationItem
    }{
        Posts:      posts,
        Pagination: generatePagination(r.URL.Path, page, len(posts)/postsLimitPerPage+1),
    }

    if r.Header.Get("HX-Request") == "true" {
        c.templates["news"].ExecuteTemplate(w, "content", data)
    } else {
        data.IsAuthenticated = isAuthenticated
        c.templates["news"].Execute(w, data)
    }
}
