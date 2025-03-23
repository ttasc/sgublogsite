package controller

import (
	"html/template"

	"github.com/go-chi/jwtauth/v5"
	"github.com/ttasc/sgublogsite/src/internal/model"
)

type Controller struct {
    model     *model.Model
    basetmpl  *template.Template
    TokenAuth *jwtauth.JWTAuth
}

func New(model *model.Model) Controller {
    return Controller{
        model:     model,
        basetmpl:  template.Must(template.New("base").ParseFiles("templates/base.tmpl")),
        TokenAuth: jwtauth.New("HS256", jwtKey, nil),
    }
}
