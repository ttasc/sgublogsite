package controller

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/ttasc/sgublogsite/src/internal/model/repos"
	"github.com/ttasc/sgublogsite/src/internal/utils"
)

const imagesFilePath = "/assets/uploads/images/"
const imagesLimitPerPage = 40

func (c *Controller) AdminImages(w http.ResponseWriter, r *http.Request) {
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }
    offset := int32 (page - 1) * imagesLimitPerPage

    images, err := c.Model.GetImages(imagesLimitPerPage, offset)
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
            "TotalPages":  int(math.Ceil(float64(total) / float64(imagesLimitPerPage))),
            "HasPrev":     page > 1,
            "HasNext":     page*imagesLimitPerPage < int(total),
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

func (c *Controller) AdminImageUpload(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(32 << 20) // 32MB max
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "File too large"})
        return
    }

    files := r.MultipartForm.File["images"]
    if len(files) == 0 {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "No files uploaded"})
    }

    for _, fileHeader := range files {
        file, err := fileHeader.Open()
        if err != nil {
            sendErrorResponse(err, w, http.StatusBadRequest,
                map[string]string{"message": "Failed to upload file (get file)"})
            return
        }
        defer file.Close()

        imageName := utils.GenerateUniqueFilename(fmt.Sprintf(".%s", imagesFilePath), fileHeader.Filename)
        imageURL := fmt.Sprintf("%s%s", imagesFilePath, imageName)

        fileHeader.Filename = imageName
        _, err = utils.SaveUploadedFile(file, fileHeader, imagesFilePath)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError,
                map[string]string{"message": "Failed to upload file (save file)"})
            return
        }

        _, err = c.Model.AddImage(fileHeader.Filename, imageURL)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError,
                map[string]string{"message": "Failed to upload file (save url)"})
            return
        }
    }
}

func (c *Controller) AdminImageBulkDelete(w http.ResponseWriter, r *http.Request) {
    var imageIDs struct {
        IDs []int `json:"ids"`
    }
    err := json.NewDecoder(r.Body).Decode(&imageIDs)
    if err != nil {
        sendErrorResponse(err, w, http.StatusBadRequest,
            map[string]string{"message": "Bad request"})
        return
    }
    for _, id := range imageIDs.IDs {
        img, _ := c.Model.GetImageByID(int32(id))
        err = utils.DeleteUploadedFile(img.Url)
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError,
                map[string]string{"message": "Failed to delete image from file system"})
        }

        err = c.Model.DeleteImage(int32(id))
        if err != nil {
            sendErrorResponse(err, w, http.StatusInternalServerError,
                map[string]string{"message": "Can't delete images"})
            return
        }
    }
    w.WriteHeader(http.StatusOK)
}
