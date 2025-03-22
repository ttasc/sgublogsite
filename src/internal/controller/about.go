package controller

import (
	"html/template"
	"net/http"
	"sgublogsite/src/internal/model"
)

func About(w http.ResponseWriter, r *http.Request) {
    about, _ := model.New().GetSiteAbout()
    data := struct {
        About string
    }{
        About: about,
    }
    tmpl, err := template.Must(basetmpl.Clone()).ParseFiles("templates/about.tmpl")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if r.Header.Get("HX-Request") == "true" {
        tmpl.ExecuteTemplate(w, "content", data)
    } else {
        tmpl.Execute(w, data)
    }
}
