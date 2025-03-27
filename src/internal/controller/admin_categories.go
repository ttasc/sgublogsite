package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (c *Controller) AdminCategories(w http.ResponseWriter, r *http.Request) {
    categories, _ := c.Model.GetCategories()
    data := struct {
        Categories []category
    }{
        Categories: buildCategoryTree(categories, 0, 0),
    }

    if r.Header.Get("HX-Request") == "true" {
        c.templates["admin_categories"].ExecuteTemplate(w, "content", data)
    } else {
        c.templates["admin_categories"].Execute(w, data)
    }
}

func (c *Controller) AdminCategoryCreate(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var category category
    err := decoder.Decode(&category)
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "Bad request"})
        return
    }
    err = c.Model.AddCategory(category.ParentID, category.Name, category.Slug)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            sendErrorResponse(err, w, http.StatusConflict,
                map[string]string{"message": "Post slug already exists"})
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (c *Controller) AdminCategoryUpdate(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var category category
    err := decoder.Decode(&category)
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "Bad request"})
        return
    }
    err = c.Model.UpdateCategory(category.ID, category.Name, category.Slug)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            sendErrorResponse(err, w, http.StatusConflict,
                map[string]string{"message": "Post slug already exists"})
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (c *Controller) AdminCategoryMove(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "Bad request"})
        return
    }
    var mvReq struct {
        NewPID *int `json:"new_parent_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&mvReq); err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "Bad request"})
        return
    }

    if mvReq.NewPID != nil && *mvReq.NewPID == id {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "Can't move to self"})
        return
    }

    if c.isCircular(id, mvReq.NewPID) {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "Can't move to circular path"})
        return
    }

    newpid := int32(0)
    if mvReq.NewPID != nil {
        newpid = int32(*mvReq.NewPID)
    }
    err = c.Model.UpdateCategoryParent(int32(id), newpid)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (c *Controller) isCircular(categoryID int, newParentID *int) bool {
    // Trường hợp chuyển thành root (không có cha)
    if newParentID == nil {
        return false
    }

    // Không cho phép làm cha của chính mình
    if *newParentID == categoryID {
        return true
    }

    // Theo dõi các parent đã kiểm tra để tránh lặp vô hạn
    checkedParents := make(map[int]bool)
    currentParentID := newParentID

    for currentParentID != nil {
        // Nếu parent hiện tại trùng với categoryID → Vòng lặp
        if *currentParentID == categoryID {
            return true
        }

        // Kiểm tra xem parent đã được kiểm tra chưa (phòng trường hợp database bị lỗi)
        if _, exists := checkedParents[*currentParentID]; exists {
            return false // Hoặc panic tùy logic
        }
        checkedParents[*currentParentID] = true

        // Lấy parent của parent hiện tại
        parentID := new(int)
        pid, err := c.Model.GetParentCategoryID(int32(*currentParentID))
        if err == nil {
            *parentID = int(pid)
        }

        if err != nil {
            return false // Xử lý lỗi tùy trường hợp
        }

        currentParentID = parentID
    }

    return false
}
