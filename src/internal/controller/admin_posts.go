package controller

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"unsafe"

	"github.com/ttasc/sgublogsite/src/internal/model/repos"
)

type showPost struct {
    repos.GetAllPostsRow
}

func (p *showPost) GetCategories() []string {
    if !p.Categories.Valid && p.Categories.String != "" {
        return []string{}
    }
    return strings.Split(p.Categories.String, ",")
}

func (c *Controller) AdminPosts(w http.ResponseWriter, r *http.Request) {
    // Handle filters
    // query := r.URL.Query()
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }

    // search := query.Get("q")
    // startDate := query.Get("start_date")
    // endDate := query.Get("end_date")
    // statusFilter := query.Get("status")
    // privateFilter := query.Get("private")

    // // Build database query
    // dbQuery := db.Model(&Post{})
    //
    // if search != "" {
    //     dbQuery = dbQuery.Where("title LIKE ?", "%"+search+"%")
    // }
    //
    // if startDate != "" {
    //     dbQuery = dbQuery.Where("created_at >= ?", startDate)
    // }
    //
    // if endDate != "" {
    //     dbQuery = dbQuery.Where("created_at <= ?", endDate)
    // }
    //
    // if statusFilter != "" && statusFilter != "all" {
    //     dbQuery = dbQuery.Where("status = ?", statusFilter)
    // }
    //
    // if privateFilter != "" && privateFilter != "all" {
    //     isPrivate, _ := strconv.ParseBool(privateFilter)
    //     dbQuery = dbQuery.Where("private = ?", isPrivate)
    // }

    // Pagination
    offset := int32 (page - 1) * postsLimitPerPage

    getPosts, err := c.Model.GetAllPosts(postsLimitPerPage, offset)
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
        c.templates["admin_posts"].ExecuteTemplate(w, "content", data)
    } else {
        c.templates["admin_posts"].Execute(w, data)
    }
}
