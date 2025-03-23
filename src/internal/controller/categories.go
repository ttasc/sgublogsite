package controller

import (
	"html/template"
	"net/http"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type category struct {
    ID       int
    Name     string
    ParentID int
    Children []category
    Level    int
}

func (c *Controller) Categories(w http.ResponseWriter, r *http.Request) {
    categories, _ := c.model.GetCategories()
    tags, _       := c.model.GetTagNames()

    data := struct {
        IsAuthenticated bool
        Categories []category
        Tags       []string
    }{
        Categories: buildCategoryTree(categories, 0, 0),
        Tags:       tags,
    }

    tmpl, _ := template.Must(c.basetmpl.Clone()).ParseFiles("templates/categories.tmpl")
    if r.Header.Get("HX-Request") == "true" {
        tmpl.ExecuteTemplate(w, "content", data)
    } else {
        _, claims, err := jwtauth.FromContext(r.Context())
        data.IsAuthenticated = (claims != nil && err == nil)
        tmpl.Execute(w, data)
    }
}

func buildCategoryTree(categories []repos.Category, parentID int, level int) []category {
    var result []category
    for _, cat := range categories {
        tmpcat := category{
            ID      : int(cat.CategoryID),
            Name    : cat.Name,
            ParentID: int(cat.ParentCategoryID.Int32),
        }
        if tmpcat.ParentID == parentID {
            tmpcat.Level = level
            tmpcat.Children = buildCategoryTree(categories, tmpcat.ID, level+1)
            result = append(result, tmpcat)
        }
    }
    return result
}

func (c *Controller) CategoryPosts(w http.ResponseWriter, r *http.Request) {
    isAuthenticated := false
    _, claims, err := jwtauth.FromContext(r.Context())
    if claims != nil && err == nil {
        isAuthenticated = true
    }

    // categoryID := strings.TrimPrefix(r.URL.Path, "/category/")
    categoryID, _ := strconv.Atoi(chi.URLParam(r, "id"))
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }
    offset := int32 (page - 1) * postsLimitPerPage

    posts, _ := c.model.GetPostsByCategoryID(
        int32(categoryID),
        postsLimitPerPage,
        offset,
        string(repos.PostsStatusPublished),
        isAuthenticated,
    )
    data := struct {
        isAuthenticated bool
        Posts       []repos.GetPostsByCategoryIDRow
        Pagination  []paginationItem
    }{
        Posts:      posts,
        Pagination: generatePagination(r.URL.Path, page, len(posts)/postsLimitPerPage+1),
    }

    tmpl, _ := template.Must(c.basetmpl.Clone()).ParseFiles("templates/news.tmpl")
    if r.Header.Get("HX-Request") == "true" {
        tmpl.ExecuteTemplate(w, "content", data)
    } else {
        data.isAuthenticated = isAuthenticated
        tmpl.Execute(w, data)
    }
}
