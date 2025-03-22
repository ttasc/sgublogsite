package controller

import (
	"html/template"
	"net/http"
	"sgublogsite/src/internal/model"
)

func Contact(w http.ResponseWriter, r *http.Request) {
    contactInfo, _ := model.New().GetContactInfo()
    data := struct {
        Address string
        Email   string
        Phone   string
    }{
        Address: contactInfo.ContactAddress.String,
        Email:   contactInfo.ContactEmail.String,
        Phone:   contactInfo.ContactPhone.String,
    }
    tmpl, err := template.Must(basetmpl.Clone()).ParseFiles("templates/contact.tmpl")
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
