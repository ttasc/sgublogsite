package controller

import (
	"html/template"
	"net/http"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

var ValidRoles = []string{
    string(repos.UsersRoleAdmin),
    // string(repos.UsersRoleAuthor),
    string(repos.UsersRoleSubscriber),
}

func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
    // id := strings.TrimPrefix(r.URL.Path, "/post/")
    id, _ := strconv.Atoi(chi.URLParam(r, "id"))

    post, err := c.Model.GetPostByID(int32(id))
    if err != nil {
        http.Error(w, "Bài viết không tồn tại", http.StatusNotFound)
        return
    }
    post.Body = string(template.HTML(post.Body))

    _, claims, err := jwtauth.FromContext(r.Context())
    isAuthenticated := (claims != nil && err == nil)
    if post.Private && !isAuthenticated {
        http.Error(w, "Bạn không có quyền truy cập", http.StatusForbidden)
        return
    }

    tags, _ := c.Model.GetTagsByPostID(int32(id))

    data := struct {
        IsAuthenticated bool
        Post repos.GetPostByIDRow
        Tags []repos.Tag
    }{
        Post: post,
        Tags: tags,
    }

    if r.Header.Get("HX-Request") == "true" {
        c.templates["post"].ExecuteTemplate(w, "content", data)
    } else {
        data.IsAuthenticated = isAuthenticated
        c.templates["post"].Execute(w, data)
    }
}
