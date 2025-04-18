package controller

import (
	"net/http"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
)


func (c *Controller) Search(w http.ResponseWriter, r *http.Request) {
    isAuthenticated := false
    _, claims, err := jwtauth.FromContext(r.Context())
    if claims != nil && err == nil {
        isAuthenticated = true
    }

    query := r.URL.Query().Get("q")
    if query == "" {
        http.Error(w, "Missing search query", http.StatusBadRequest)
        return
    }
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }
    offset := int32 (page - 1) * postsLimitPerPage

    posts, _ := c.Model.SearchPosts(
        query,
        postsLimitPerPage,
        offset,
        string(repos.PostsStatusPublished),
        isAuthenticated,
    )
    data := struct {
        IsAuthenticated bool
        SearchQuery string
        Posts       []repos.FindPostsRow
        Pagination  []paginationItem
    }{
        SearchQuery: query,
        Posts:       posts,
        Pagination:  generatePagination(r.URL.Path+"?q="+query, page, len(posts)/postsLimitPerPage+1),
    }
    if r.Header.Get("HX-Request") == "true" {
        c.templates["search"].ExecuteTemplate(w, "content", data)
    } else {
        data.IsAuthenticated = isAuthenticated
        c.templates["search"].Execute(w, data)
    }
}
