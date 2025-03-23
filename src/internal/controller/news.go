package controller

import (
	"html/template"
	"net/http"
	"github.com/ttasc/sgublogsite/src/internal/model"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
)

func News(w http.ResponseWriter, r *http.Request) {
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

    posts, _ := model.New().GetPostsByCategorySlug(
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

    tmpl, _ := template.Must(basetmpl.Clone()).ParseFiles("templates/news.tmpl")
    if r.Header.Get("HX-Request") == "true" {
        tmpl.ExecuteTemplate(w, "content", data)
    } else {
        data.IsAuthenticated = isAuthenticated
        tmpl.Execute(w, data)
    }
}
