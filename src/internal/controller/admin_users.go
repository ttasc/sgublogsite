package controller

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
)

func (c *Controller) AdminUsers(w http.ResponseWriter, r *http.Request) {
    users, _ := c.Model.GetUsers()
    data := struct {
        Users []repos.GetAllUsersRow
    }{
        Users: users,
    }

    if r.Header.Get("HX-Request") == "true" {
        c.templates["admin_users"].ExecuteTemplate(w, "content", data)
    } else {
        c.templates["admin_users"].Execute(w, data)
    }
}

func (c *Controller) AdminUser(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(chi.URLParam(r, "id"))
    user, _ := c.Model.GetUserByID(int32(id))
    if r.Header.Get("HX-Request") == "true" {
        c.templates["admin_user"].ExecuteTemplate(w, "content", user)
    } else {
        c.templates["admin_user"].Execute(w, user)
    }
}

func (c *Controller) AdminUserDelete(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(chi.URLParam(r, "id"))
    err := c.Model.DeleteUser(int32(id))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}
