package controller

import (
	"net/http"
)

func (c *Controller) AdminWelcome(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") == "true" {
        c.templates["admin_welcome"].ExecuteTemplate(w, "content", nil)
    } else {
        c.templates["admin_welcome"].Execute(w, nil)
    }
}
