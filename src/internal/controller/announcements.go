package controller

import (
	"html/template"
	"net/http"
	"sgublogsite/src/internal/model"
	"sgublogsite/src/internal/model/repos"
	"strconv"
	"strings"

	"github.com/go-chi/jwtauth/v5"
)

type monthGroup struct {
    MonthYear string
    Items     []repos.GetPostsByCategorySlugRow
}

func Announcements(w http.ResponseWriter, r *http.Request) {
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

    tmpl, _ := template.Must(basetmpl.Clone()).ParseFiles("templates/announcements.tmpl")
    if r.Header.Get("HX-Request") == "true" {
        tmpl.ExecuteTemplate(w, "content", data)
    } else {
        data.IsAuthenticated = (claims != nil && err == nil)
        tmpl.Execute(w, data)
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
