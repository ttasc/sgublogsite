package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
	"github.com/ttasc/sgublogsite/src/internal/utils"
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

func (c *Controller) AdminUserNew(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") == "true" {
        c.templates["admin_user_new"].ExecuteTemplate(w, "content", nil)
    } else {
        c.templates["admin_user_new"].Execute(w, nil)
    }
}

func (c *Controller) AdminUserCreate(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(32 << 20) // 32MB max
    if err != nil {
        sendErrorResponse(w, http.StatusBadRequest, map[string]string{"message": "File too large"})
        return
    }

    firstname := r.FormValue("first_name")
    lastname := r.FormValue("last_name")
    phone := r.FormValue("phone")
    email := r.FormValue("email")
    password, err := utils.HashPassword(r.FormValue("password"))
    if err != nil {
        sendErrorResponse(w, http.StatusInternalServerError, map[string]string{"message": "Server internal error"})
        return
    }
    role := repos.NullUsersRole{UsersRole: repos.UsersRole(r.FormValue("role")), Valid: true}

    user := repos.User{
        Firstname: firstname,
        Lastname:  lastname,
        Phone:     phone,
        Email:     email,
        Password:  password,
        Role:      role,
    }

    userID, err := c.Model.AddUser(user)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            sendErrorResponse(w, http.StatusConflict, map[string]string{"message": "Email or phone already exists"})
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    avatarURL := "/assets/avatar.def.png"
    file, handler, err := r.FormFile("avatar")
    if err == nil {
        defer file.Close()
        filename := handler.Filename
        fileExt := filename[strings.LastIndex(filename, ".")+1:]
        filename = strconv.Itoa(int(userID)) + "." + fileExt
        handler.Filename = filename
        avatarURL, err = utils.SaveUploadedFile(file, handler)
        if err != nil {
            sendErrorResponse(w, http.StatusInternalServerError, map[string]string{"message": "Failed to upload file (save file)"})
            return
        }

        imgID, err := c.Model.AddImage(repos.Image{
            Url:          avatarURL,
            Name:         sql.NullString{String: handler.Filename, Valid: true},
        })
        if err != nil {
            sendErrorResponse(w, http.StatusInternalServerError, map[string]string{"message": "Failed to upload file (save url)"})
            return
        }

        c.Model.UpdateUserInfo(repos.User{
            UserID:         userID,
            Firstname:      firstname,
            Lastname:       lastname,
            Phone:          phone,
            Email:          email,
            ProfilePicID:   sql.NullInt32{Int32: imgID, Valid: true},
        })
    }

    http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
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

func sendErrorResponse(w http.ResponseWriter, statusCode int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}
