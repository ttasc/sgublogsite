package controller

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unsafe"

	"github.com/ttasc/sgublogsite/src/internal/model/repos"
)

type showPost struct {
    repos.GetFilteredPostsRow
}

func (p *showPost) GetCategories() []string {
    if !p.Categories.Valid && p.Categories.String != "" {
        return []string{}
    }
    return strings.Split(p.Categories.String, ",")
}

func (c *Controller) AdminPosts(w http.ResponseWriter, r *http.Request) {
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }

    query := r.URL.Query()
    search := query.Get("q")
    statusFilter := query.Get("status")
    privateFilter := query.Get("private")

    var getStatus repos.PostsStatus
    var isPrivate bool
    switch {
    case statusFilter != "" && statusFilter != "all":
        getStatus = repos.PostsStatus(statusFilter)
    case privateFilter != "" && privateFilter != "all":
        isPrivate, _ = strconv.ParseBool(privateFilter)
    }

    // Pagination
    offset := int32 (page - 1) * postsLimitPerPage

    getPosts, err := c.Model.GetFilteredPosts(postsLimitPerPage, offset, search, getStatus, isPrivate)
    posts := *(*[]showPost)(unsafe.Pointer(&getPosts))
    if err != nil {
        sendErrorResponse(err, w, http.StatusInternalServerError,
            map[string]string{"message": "Can't get posts"})
        return
    }

    total, err := c.Model.CountPosts()
    if err != nil {
        sendErrorResponse(err, w, http.StatusInternalServerError,
            map[string]string{"message": "Can't count posts"})
        return
    }
    data := map[string]any{
        "Posts": posts,
        "Pagination": map[string]any{
            "CurrentPage": page,
            "TotalPages":  int(math.Ceil(float64(total) / float64(postsLimitPerPage))),
            "HasPrev":     page > 1,
            "HasNext":     page*postsLimitPerPage < int(total),
            "NextPage":    page + 1,
            "PrevPage":    page - 1,
        },
    }

    if r.Header.Get("HX-Request") == "true" {
        if r.Header.Get("HX-Target") == "postTable" {
            c.templates["admin_posts"].ExecuteTemplate(w, "posts_table", data)
        } else {
            c.templates["admin_posts"].ExecuteTemplate(w, "content", data)
        }
    } else {
        c.templates["admin_posts"].Execute(w, data)
    }
}

func (c *Controller) AdminPostBulkDelete(w http.ResponseWriter, r *http.Request) {
    var postIDs struct {
        IDs []int `json:"ids"`
    }
    err := json.NewDecoder(r.Body).Decode(&postIDs)
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "Bad request"})
        return
    }
    for _, id := range postIDs.IDs {
        err = c.Model.DeletePost(int32(id))
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError,
                map[string]string{"message": "Can't delete selected posts"})
            return
        }
    }
    w.WriteHeader(http.StatusOK)
}
