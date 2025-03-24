package controller

import (
	"net/http"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
	"strconv"
	"strings"

	"github.com/go-chi/jwtauth/v5"
)

type monthGroup struct {
    MonthYear string
    Items     []repos.GetPostsByCategorySlugRow
}

func (c *Controller) Announcements(w http.ResponseWriter, r *http.Request) {
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
        "announcements",
        postsLimitPerPage,
        offset,
        string(repos.PostsStatusPublished),
        isAuthenticated,
    )
    data := struct {
        IsAuthenticated      bool
        GroupedAnnouncements []monthGroup
        Pagination           []paginationItem
    }{
        GroupedAnnouncements: groupAnnouncementsByMonth(posts),
        Pagination:           generatePagination(r.URL.Path, page, len(posts)/postsLimitPerPage+1),
    }

    if r.Header.Get("HX-Request") == "true" {
        c.templates["announcements"].ExecuteTemplate(w, "content", data)
    } else {
        data.IsAuthenticated = (claims != nil && err == nil)
        c.templates["announcements"].Execute(w, data)
    }
}

func groupAnnouncementsByMonth(announcements []repos.GetPostsByCategorySlugRow) []monthGroup {
    groups := make(map[string][]repos.GetPostsByCategorySlugRow)

    for _, a := range announcements {
        monthYear := a.CreatedAt.Format("01/2006")
        groups[monthYear] = append(groups[monthYear], a)
    }

    var result []monthGroup
    for key, val := range groups {
        result = append(result, monthGroup{
            MonthYear: translateMonth(key),
            Items:     val,
        })
    }

    return result
}

func translateMonth(monthYear string) string {
    parts := strings.Split(monthYear, "/")
    return "Tháng " + parts[0] + " năm " + parts[1]
}
