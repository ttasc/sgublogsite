package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ttasc/sgublogsite/src/internal/model/repos"
)

func (c *Controller) AdminTags(w http.ResponseWriter, r *http.Request) {
    tags, _ := c.Model.GetTags()
    data := struct {
        Tags []repos.Tag
    }{
        Tags: tags,
    }
    if r.Header.Get("HX-Request") == "true" {
        c.templates["admin_tags"].ExecuteTemplate(w, "content", data)
    } else {
        c.templates["admin_tags"].Execute(w, data)
    }
}

func (c *Controller) AdminTagCreate(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var tag repos.Tag
    err := decoder.Decode(&tag)
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "Bad request"})
        return
    }
    err = c.Model.AddTag(tag.Name, tag.Slug)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            sendErrorResponse(err, w, http.StatusConflict,
                map[string]string{"message": "Tag slug already exists"})
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (c *Controller) AdminTagUpdate(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var tag repos.Tag
    err := decoder.Decode(&tag)
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "Bad request"})
        return
    }
    err = c.Model.UpdateTag(tag.TagID, tag.Name, tag.Slug)
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

func (c *Controller) AdminTagDelete(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "Bad request"})
        return
    }
    err = c.Model.DeleteTag(int32(id))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (c *Controller) AdminTagBulkDelete(w http.ResponseWriter, r *http.Request) {
    var tagIDs struct {
        IDs []int `json:"ids"`
    }
    err := json.NewDecoder(r.Body).Decode(&tagIDs)
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "Bad request"})
        return
    }
    for _, id := range tagIDs.IDs {
        err = c.Model.DeleteTag(int32(id))
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError,
                map[string]string{"message": "Can't delete tag"})
            return
        }
    }
    w.WriteHeader(http.StatusOK)
}
