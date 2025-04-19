package controller

import (
	"html/template"
	"net/http"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type htmlpost struct {
    repos.GetPostByIDRow
    HTMLBody template.HTML
}

func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
    // id := strings.TrimPrefix(r.URL.Path, "/post/")
    id, _ := strconv.Atoi(chi.URLParam(r, "id"))

    rawPost, err := c.Model.GetPostByID(int32(id))
    if err != nil {
        http.Error(w, "Bài viết không tồn tại", http.StatusNotFound)
        return
    }
    post := htmlpost{GetPostByIDRow: rawPost, HTMLBody: template.HTML(rawPost.Body)}

    _, claims, err := jwtauth.FromContext(r.Context())
    isAuthenticated := (claims != nil && err == nil)
    if post.Private && !isAuthenticated {
        http.Error(w, "Bạn không có quyền truy cập", http.StatusForbidden)
        return
    }

    tags, _ := c.Model.GetTagsByPostID(int32(id))

    data := struct {
        IsAuthenticated bool
        Post htmlpost
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
