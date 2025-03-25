package controller

import (
	"net/http"
)

func (c *Controller) AdminDashboard(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") == "true" {
        c.templates["admin_dashboard"].ExecuteTemplate(w, "content", nil)
    } else {
        c.templates["admin_dashboard"].Execute(w, nil)
    }
}
