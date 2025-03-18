package server

import (
	"io"
	"slices"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

var t = &Template{
    templates: template.Must(
        template.New("").Funcs(
            template.FuncMap{
                "hasRole": func(roles []string, target string) bool {
                    if slices.Contains(roles, target) {
                            return true
                    }
                    return false
                },
            },
        ).ParseGlob("src/statics/templates/*.html"),
    ),
}
