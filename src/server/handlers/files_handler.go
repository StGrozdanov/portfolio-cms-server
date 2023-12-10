package handlers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"portfolio-cms-server/internal/files"
	"portfolio-cms-server/utils"
)

func UploadCV(ginCtx *gin.Context) {
	file, err := ginCtx.FormFile("file")
	if err != nil {
		ginCtx.JSON(
			http.StatusBadRequest,
			map[string]interface{}{"message": "the provided file should be with the name 'file'."},
		)
		return
	}

	if file.Header.Get("Content-Type") != "application/pdf" {
		ginCtx.JSON(
			http.StatusBadRequest,
			map[string]interface{}{"message": "provided file can only be of type pdf"},
		)
		return
	}

	cvLink, err := files.UploadCV(file)
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on attempting to upload a CV")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusCreated, map[string]interface{}{"cvLink": cvLink})
}

func UploadProjectImage(ginCtx *gin.Context) {
	image, _ := ginCtx.FormFile("image")

	projectTitle, found := ginCtx.GetPostForm("projectTitle")
	if !found {
		ginCtx.JSON(
			http.StatusBadRequest,
			map[string]interface{}{"error": "invalid parameters, expected projectTitle"},
		)
		return
	}

	projectImages, err := files.UploadProjectImage(image, projectTitle)
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on attempting to upload a project image")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusCreated, map[string]interface{}{"project_images": projectImages})
}

func UploadJobImage(ginCtx *gin.Context) {
	image, _ := ginCtx.FormFile("image")

	company, found := ginCtx.GetPostForm("companyName")
	if !found {
		ginCtx.JSON(
			http.StatusBadRequest,
			map[string]interface{}{"error": "invalid parameters, expected companyName"},
		)
		return
	}

	jobImages, err := files.UploadJobImage(image, company)
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on attempting to upload a job image")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusCreated, map[string]interface{}{"job_images": jobImages})
}
