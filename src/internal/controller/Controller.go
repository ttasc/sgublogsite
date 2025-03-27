package controller

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/ttasc/sgublogsite/src/internal/model"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
)

var ValidRoles = []string{
    string(repos.UsersRoleAdmin),
    // string(repos.UsersRoleAuthor),
    string(repos.UsersRoleSubscriber),
}

type Controller struct {
    templates map[string]*template.Template
    Model     *model.Model
    TokenAuth *jwtauth.JWTAuth
}

func New(model *model.Model) Controller {
    return Controller{
        templates: initTemplates(),
        Model:     model,
        TokenAuth: jwtauth.New("HS256", jwtKey, nil),
    }
}

func initTemplates() map[string]*template.Template {
    base  := template.Must(template.New("base").ParseFiles("templates/base.tmpl"))
    admin := template.Must(template.New("admin").ParseFiles("templates/admin/admin.tmpl"))
    return map[string]*template.Template{
        "profile":          template.Must(template.ParseFiles("templates/profile.tmpl")),

        "base":             base,
        "announcements":    parseBase(base, "templates/announcements.tmpl"),
        "post":             parseBase(base, "templates/post.tmpl"),
        "search":           parseBase(base, "templates/search.tmpl"),
        "categories":       parseBase(base, "templates/categories.tmpl"),
        "news":             parseBase(base, "templates/news.tmpl"),
        "contact":          parseBase(base, "templates/contact.tmpl"),
        "about":            parseBase(base, "templates/about.tmpl"),
        "home":             parseBase(base, "templates/home.tmpl"),

        "admin":            admin,
        "admin_dashboard":  parseBase(admin, "templates/admin/dashboard.tmpl"),
        "admin_users":      parseBase(admin, "templates/admin/users.tmpl"),
        "admin_user":       parseBase(admin, "templates/admin/user.tmpl"),
        "admin_user_new":   parseBase(admin, "templates/admin/user_new.tmpl"),
        "admin_categories": parseBase(admin, "templates/admin/categories.tmpl"),
    }
}

func parseBase(base *template.Template, filename string) *template.Template {
    return template.Must(template.Must(base.Clone()).ParseFiles(filename))
}

func sendErrorResponse(err error, w http.ResponseWriter, statusCode int, data any) {
    log.Println(err)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}
