package controller

import (
	"net/http"
	"strconv"
	"unsafe"

	"github.com/ttasc/sgublogsite/src/internal/model"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type category struct {
    ID       int32      `json:"id"`
    Name     string     `json:"name"`
    Slug     string     `json:"slug"`
    ParentID int32      `json:"parent_id"`
    Children []category `json:"children"`
    Level    int        `json:"level"`
}

func (c *Controller) Categories(w http.ResponseWriter, r *http.Request) {
    categories, _ := c.Model.GetCategories()
    tags, _       := c.Model.GetTags()

    data := struct {
        IsAuthenticated bool
        Categories []category
        Tags       []repos.Tag
    }{
        Categories: buildCategoryTree(categories, 0, 0),
        Tags:       tags,
    }

    if r.Header.Get("HX-Request") == "true" {
        c.templates["categories"].ExecuteTemplate(w, "content", data)
    } else {
        _, claims, err := jwtauth.FromContext(r.Context())
        data.IsAuthenticated = (claims != nil && err == nil)
        c.templates["categories"].Execute(w, data)
    }
}

func buildCategoryTree(categories []repos.Category, parentID int32, level int) []category {
    var result []category
    for _, cat := range categories {
        tmpcat := category{
            ID      : cat.CategoryID,
            Name    : cat.Name,
            Slug    : cat.Slug,
            ParentID: cat.ParentCategoryID.Int32,
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

    var posts []model.GetPost
    if categoryID == -1 {
        getPosts, _ := c.Model.GetUncategorizedPosts(
            postsLimitPerPage,
            offset,
            string(repos.PostsStatusPublished),
            isAuthenticated,
        )
        posts = *(*[]model.GetPost)(unsafe.Pointer(&getPosts))
    } else {
        getPosts, _ := c.Model.GetPostsByCategoryID(
            int32(categoryID),
            postsLimitPerPage,
            offset,
            string(repos.PostsStatusPublished),
            isAuthenticated,
        )
        posts = *(*[]model.GetPost)(unsafe.Pointer(&getPosts))
    }
    data := struct {
        isAuthenticated bool
        Posts       []model.GetPost
        Pagination  []paginationItem
    }{
        Posts:      posts,
        Pagination: generatePagination(r.URL.Path, page, len(posts)/postsLimitPerPage+1),
    }

    if r.Header.Get("HX-Request") == "true" {
        c.templates["news"].ExecuteTemplate(w, "content", data)
    } else {
        data.isAuthenticated = isAuthenticated
        c.templates["news"].Execute(w, data)
    }
}

func (c *Controller) TagPosts(w http.ResponseWriter, r *http.Request) {
    isAuthenticated := false
    _, claims, err := jwtauth.FromContext(r.Context())
    if claims != nil && err == nil {
        isAuthenticated = true
    }

    // tagID := strings.TrimPrefix(r.URL.Path, "/tag/")
    tagID, _ := strconv.Atoi(chi.URLParam(r, "id"))
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }
    offset := int32 (page - 1) * postsLimitPerPage

    posts, _ := c.Model.GetPostsByTagID(
        int32(tagID),
        postsLimitPerPage,
        offset,
        string(repos.PostsStatusPublished),
        isAuthenticated,
    )
    data := struct {
        isAuthenticated bool
        Posts       []repos.GetPostsByTagIDRow
        Pagination  []paginationItem
    }{
        Posts:      posts,
        Pagination: generatePagination(r.URL.Path, page, len(posts)/postsLimitPerPage+1),
    }

    if r.Header.Get("HX-Request") == "true" {
        c.templates["news"].ExecuteTemplate(w, "content", data)
    } else {
        data.isAuthenticated = isAuthenticated
        c.templates["news"].Execute(w, data)
    }
}
