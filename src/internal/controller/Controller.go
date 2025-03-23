package controller

import (
	"html/template"

	"github.com/go-chi/jwtauth/v5"
	"github.com/ttasc/sgublogsite/src/internal/model"
)

type Controller struct {
    basetmpl  *template.Template
    Model     *model.Model
    TokenAuth *jwtauth.JWTAuth
}

func New(model *model.Model) Controller {
    return Controller{
        basetmpl:  template.Must(template.New("base").ParseFiles("templates/base.tmpl")),
        Model:     model,
        TokenAuth: jwtauth.New("HS256", jwtKey, nil),
    }
}
