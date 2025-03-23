package controller

import (
	"html/template"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func (c *Controller) Contact(w http.ResponseWriter, r *http.Request) {
    contactInfo, _ := c.Model.GetContactInfo()
    data := struct {
        IsAuthenticated bool
        Address string
        Email   string
        Phone   string
    }{
        Address: contactInfo.ContactAddress.String,
        Email:   contactInfo.ContactEmail.String,
        Phone:   contactInfo.ContactPhone.String,
    }
    tmpl, err := template.Must(c.basetmpl.Clone()).ParseFiles("templates/contact.tmpl")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if r.Header.Get("HX-Request") == "true" {
        tmpl.ExecuteTemplate(w, "content", data)
    } else {
        _, claims, err := jwtauth.FromContext(r.Context())
        data.IsAuthenticated = (claims != nil && err == nil)
        tmpl.Execute(w, data)
    }
}
