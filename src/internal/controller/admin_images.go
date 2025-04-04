package controller

import (
	"math"
	"net/http"
	"strconv"

	"github.com/ttasc/sgublogsite/src/internal/model/repos"
)

func (c *Controller) AdminImages(w http.ResponseWriter, r *http.Request) {
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }
    offset := int32 (page - 1) * postsLimitPerPage

    images, err := c.Model.GetImages(postsLimitPerPage, offset)
    if err != nil {
        sendErrorResponse(err, w, http.StatusInternalServerError,
            map[string]string{"message": "Can't get images"})
        return
    }

    total, err := c.Model.CountImages()
    if err != nil {
        sendErrorResponse(err, w, http.StatusInternalServerError,
            map[string]string{"message": "Can't count images"})
        return
    }

    data := struct {
        Images []repos.Image
        Pagination map[string]any
    } {
        Images: images,
        Pagination: map[string]any{
            "CurrentPage": page,
            "TotalPages":  int(math.Ceil(float64(total) / float64(postsLimitPerPage))),
            "HasPrev":     page > 1,
            "HasNext":     page*postsLimitPerPage < int(total),
            "NextPage":    page + 1,
            "PrevPage":    page - 1,
        },
    }

    if r.Header.Get("HX-Request") == "true" {
        c.templates["admin_images"].ExecuteTemplate(w, "content", data)
    } else {
        c.templates["admin_images"].Execute(w, data)
    }
}
