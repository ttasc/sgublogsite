package controller

import (
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
    if r.Header.Get("HX-Request") == "true" {
        c.templates["contact"].ExecuteTemplate(w, "content", data)
    } else {
        _, claims, err := jwtauth.FromContext(r.Context())
        data.IsAuthenticated = (claims != nil && err == nil)
        c.templates["contact"].Execute(w, data)
    }
}
