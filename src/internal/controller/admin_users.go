package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
	"github.com/ttasc/sgublogsite/src/internal/utils"
)

const avatarsFilePath = "/assets/uploads/avatars/"

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
    data := struct {
        User repos.GetUserByIDRow
        Roles []string
    } {
        User: user,
        Roles: ValidRoles,
    }
    if r.Header.Get("HX-Request") == "true" {
        c.templates["admin_user"].ExecuteTemplate(w, "content", data)
    } else {
        c.templates["admin_user"].Execute(w, data)
    }
}

func (c *Controller) AdminUserNew(w http.ResponseWriter, r *http.Request) {
    data := struct {
        Roles []string
    } {
        Roles: ValidRoles,
    }
    if r.Header.Get("HX-Request") == "true" {
        c.templates["admin_user_new"].ExecuteTemplate(w, "content", data)
    } else {
        c.templates["admin_user_new"].Execute(w, data)
    }
}

func (c *Controller) AdminUserCreate(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(32 << 20) // 32MB max
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "File too large"})
        return
    }

    firstname := r.FormValue("first_name")
    lastname := r.FormValue("last_name")
    phone := r.FormValue("phone")
    email := r.FormValue("email")
    password, err := utils.HashPassword(r.FormValue("password"))
    if err != nil {
        sendErrorResponse(err, w, http.StatusInternalServerError,
            map[string]string{"message": "Server internal error"})
        return
    }
    role := repos.UsersRole(r.FormValue("role"))

    userID, err := c.Model.AddUser(firstname, lastname, phone, email, password, role)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            sendErrorResponse(err, w, http.StatusConflict,
                map[string]string{"message": "Email or phone already exists"})
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
        handler.Filename = fmt.Sprintf("avatar-user%d.%s", userID, fileExt)
        avatarURL, err := utils.SaveUploadedFile(file, handler, avatarsFilePath)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError,
                map[string]string{"message": "Failed to upload file (save file)"})
            return
        }

        imgID, err := c.Model.AddImage(handler.Filename, avatarURL)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError,
                map[string]string{"message": "Failed to upload file (save url)"})
            return
        }

        c.Model.UpdateUserAvatarID(userID, imgID)
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
        sendErrorResponse(err, w, http.StatusNotFound,
            map[string]string{"message": "User not found"})
        return
    }
    if user.UserID == currentUser.UserID && user.Role == repos.UsersRoleAdmin {
        sendErrorResponse(err, w, http.StatusForbidden,
            map[string]string{"message": "You can't delete yourself"})
        return
    }
    if user.AvatarID.Valid {
        err = utils.DeleteUploadedFile(user.Avatar.String)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError,
                map[string]string{"message": "Failed to delete image from file system"})
        }
        err = c.Model.DeleteImage(user.AvatarID.Int32)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError,
                map[string]string{"message": "Failed to delete image from database"})
        }
    }
    err = c.Model.DeleteUser(int32(id))
    if err != nil {
        sendErrorResponse(err, w, http.StatusInternalServerError,
            map[string]string{"message": "Failed to delete user"})
        return
    }

    http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}


func (c *Controller) updateUser(userID int32, w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(32 << 20) // 32MB max
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "File too large"})
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
            sendErrorResponse(err, w, http.StatusInternalServerError,
                map[string]string{"message": "Server internal error"})
            return
        }
    }
    role := repos.UsersRole(r.FormValue("role"))

    imgID, _ := c.Model.GetUserAvatarID(userID)

    file, handler, err := r.FormFile("avatar")
    if err == nil {
        defer file.Close()
        filename := handler.Filename
        fileExt := filename[strings.LastIndex(filename, ".")+1:]
        handler.Filename = fmt.Sprintf("avatar-user%d.%s", userID, fileExt)
        avatarURL, err := utils.SaveUploadedFile(file, handler, avatarsFilePath)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError,
                map[string]string{"message": "Failed to upload file (save file)"})
            return
        }

        if imgID != 0 {
            err = c.Model.UpdateImageURL(imgID, avatarURL)
            if err != nil {
                sendErrorResponse(err, w, http.StatusInternalServerError,
                    map[string]string{"message": "Failed to upload file (save url)"})
                return
            }
        } else {
            imgID, err = c.Model.AddImage(handler.Filename, avatarURL)
        }
    }

    err = c.Model.UpdateUserInfo(userID, firstname, lastname, phone, email)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            sendErrorResponse(err, w, http.StatusConflict,
                map[string]string{"message": "Email or phone already exists"})
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    c.Model.UpdateUserAvatarID(userID, imgID)
    c.Model.UpdateUserRole(userID, role)
    if password != "" {
        c.Model.UpdateUserPassword(userID, password)
    }
}
