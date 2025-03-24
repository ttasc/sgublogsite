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
    base := template.Must(template.New("base").ParseFiles("templates/base.tmpl"))
    return map[string]*template.Template{
        "profile":          template.Must(template.ParseFiles("templates/profile.tmpl")),
        "base":             base,
        "announcements":    template.Must(template.Must(base.Clone()).ParseFiles("templates/announcements.tmpl")),
        "post":             template.Must(template.Must(base.Clone()).ParseFiles("templates/post.tmpl")),
        "search":           template.Must(template.Must(base.Clone()).ParseFiles("templates/search.tmpl")),
        "categories":       template.Must(template.Must(base.Clone()).ParseFiles("templates/categories.tmpl")),
        "news":             template.Must(template.Must(base.Clone()).ParseFiles("templates/news.tmpl")),
        "contact":          template.Must(template.Must(base.Clone()).ParseFiles("templates/contact.tmpl")),
        "about":            template.Must(template.Must(base.Clone()).ParseFiles("templates/about.tmpl")),
        "home":             template.Must(template.Must(base.Clone()).ParseFiles("templates/home.tmpl")),
    }
}
