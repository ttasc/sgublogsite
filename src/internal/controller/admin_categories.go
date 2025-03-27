package controller

import (
	"encoding/json"
	"net/http"
	"strings"
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
