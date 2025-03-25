package controller

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
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
        sendErrorResponse(err, w, http.StatusBadRequest, map[string]string{"message": "File too large"})
        return
    }

    firstname := r.FormValue("first_name")
    lastname := r.FormValue("last_name")
    phone := r.FormValue("phone")
    email := r.FormValue("email")
    password, err := utils.HashPassword(r.FormValue("password"))
    if err != nil {
        sendErrorResponse(err, w, http.StatusInternalServerError, map[string]string{"message": "Server internal error"})
        return
    }
    role := repos.UsersRole(r.FormValue("role"))

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
            sendErrorResponse(err, w, http.StatusConflict, map[string]string{"message": "Email or phone already exists"})
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    file, handler, err := r.FormFile("avatar")
    if err == nil {
        defer file.Close()
        filename := handler.Filename
        fileExt := filename[strings.LastIndex(filename, ".")+1:]
        filename = strconv.Itoa(int(userID)) + "." + fileExt
        handler.Filename = filename
        avatarURL, err := utils.SaveUploadedFile(file, handler)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError, map[string]string{"message": "Failed to upload file (save file)"})
            return
        }

        imgID, err := c.Model.AddImage(repos.Image{
            Url:          avatarURL,
            Name:         sql.NullString{String: "avatar", Valid: true},
        })
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError, map[string]string{"message": "Failed to upload file (save url)"})
            return
        }

        c.Model.UpdateUserProfilePicID(userID, imgID)
    }

    http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (c *Controller) AdminUserUpdate(w http.ResponseWriter, r *http.Request) {
    userID, _ := strconv.Atoi(chi.URLParam(r, "id"))

    c.updateUser(int32(userID), w, r)

    http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (c *Controller) AdminUserDelete(w http.ResponseWriter, r *http.Request) {
    _, claims, err := jwtauth.FromContext(r.Context())
    if claims == nil || err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    currentUser, _ := c.Model.GetUserByID(int32(claims["ID"].(float64)))

    id, _ := strconv.Atoi(chi.URLParam(r, "id"))
    user, err := c.Model.GetUserByID(int32(id))
    if err != nil {
        sendErrorResponse(err, w, http.StatusNotFound, map[string]string{"message": "User not found"})
        return
    }
    if user.Role == currentUser.Role && currentUser.Role == repos.UsersRoleAdmin {
        sendErrorResponse(err, w, http.StatusForbidden, map[string]string{"message": "You can't delete yourself"})
        return
    }
    if user.ProfilePicID.Valid {
        err = utils.DeleteUploadedFile(user.ProfilePic.String)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError, map[string]string{"message": "Failed to delete image from file system"})
        }
        err = c.Model.DeleteImage(user.ProfilePicID.Int32)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError, map[string]string{"message": "Failed to delete image from database"})
        }
    }
    err = c.Model.DeleteUser(int32(id))
    if err != nil {
        sendErrorResponse(err, w, http.StatusInternalServerError, map[string]string{"message": "Failed to delete user"})
        return
    }

    http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}


func (c *Controller) updateUser(userID int32, w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(32 << 20) // 32MB max
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest, map[string]string{"message": "File too large"})
        return
    }

    firstname := r.FormValue("first_name")
    lastname := r.FormValue("last_name")
    phone := r.FormValue("phone")
    email := r.FormValue("email")
    password := r.FormValue("password")
    if password != "" {
        password, err = utils.HashPassword(password)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError, map[string]string{"message": "Server internal error"})
            return
        }
    }
    role := repos.UsersRole(r.FormValue("role"))

    imgID, err := c.Model.GetUserProfilePicID(userID)
    if err != nil {
        imgID.Valid = false
    }

    file, handler, err := r.FormFile("avatar")
    if err == nil {
        defer file.Close()
        filename := handler.Filename
        fileExt := filename[strings.LastIndex(filename, ".")+1:]
        filename = strconv.Itoa(int(userID)) + "." + fileExt
        handler.Filename = filename
        avatarURL, err := utils.SaveUploadedFile(file, handler)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError, map[string]string{"message": "Failed to upload file (save file)"})
            return
        }

        if imgID.Valid {
            err = c.Model.UpdateImageURL(imgID.Int32, avatarURL)
            if err != nil {
                sendErrorResponse(err, w, http.StatusInternalServerError, map[string]string{"message": "Failed to upload file (save url)"})
                return
            }
        } else {
            imgID.Int32, err = c.Model.AddImage(repos.Image{
                Url:          avatarURL,
                Name:         sql.NullString{String: "avatar", Valid: true},
            })
        }
    }

    err = c.Model.UpdateUserInfo(repos.User{
        UserID:         userID,
        Firstname:      firstname,
        Lastname:       lastname,
        Phone:          phone,
        Email:          email,
    })
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            sendErrorResponse(err, w, http.StatusConflict, map[string]string{"message": "Email or phone already exists"})
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    c.Model.UpdateUserProfilePicID(userID, imgID.Int32)
    c.Model.UpdateUserRole(userID, role)
    if password != "" {
        c.Model.UpdateUserPassword(userID, password)
    }
}
