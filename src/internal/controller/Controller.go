package controller

import (
	"html/template"

	"github.com/go-chi/jwtauth/v5"
	"github.com/ttasc/sgublogsite/src/internal/model"
)

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
        "admin_welcome":    parseBase(admin, "templates/admin/welcome.tmpl"),
        "admin_users":      parseBase(admin, "templates/admin/users.tmpl"),
        "admin_user":       parseBase(admin, "templates/admin/user.tmpl"),
    }
}

func parseBase(base *template.Template, filename string) *template.Template {
    return template.Must(template.Must(base.Clone()).ParseFiles(filename))
}
