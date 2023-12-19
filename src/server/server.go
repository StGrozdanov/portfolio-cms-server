package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
	"portfolio-cms-server/server/handlers"
	"portfolio-cms-server/server/middlewares"
	"portfolio-cms-server/utils"
)

var geoKey string

func setupRouter(db *geoip2.Reader) (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	router.Use(middlewares.Logger(utils.GetLogger()), gin.Recovery())
	router.Use(middlewares.CORS())

	router.GET("/healths", handlers.HealthCheck)
	router.GET("/metrics", handlers.Metrics)
	router.GET("/users/basic-info", handlers.GetBasicInfo)
	router.PUT("/users/basic-info", middlewares.AuthMiddleware(), handlers.UpdateBasicInfo)
	router.GET("/users/skills", handlers.GetSkills)
	router.PUT("/users/skills", middlewares.AuthMiddleware(), handlers.UpdateSkills)
	router.GET("/users/jobs-and-projects", handlers.GetJobsAndProjects)
	router.PUT("/users/jobs-and-projects", middlewares.AuthMiddleware(), handlers.UpdateJobsAndProjects)
	router.GET("/users/socials", handlers.GetSocials)
	router.PUT("/users/socials", middlewares.AuthMiddleware(), handlers.UpdateSocials)
	router.POST("/auth/login", handlers.Login)
	router.POST("/analytics/track", func(ginCtx *gin.Context) {
		handlers.Track(ginCtx, db)
	})

	fileAuthGroup := router.Group("/files")
	fileAuthGroup.Use(middlewares.AuthMiddleware())
	{
		fileAuthGroup.POST("/cv", handlers.UploadCV)
		fileAuthGroup.DELETE("/image", handlers.DeleteImage)

		imageGroup := fileAuthGroup.Group("")
		imageGroup.Use(middlewares.ImageContentTypeMiddleware())
		{
			imageGroup.POST("/project-image", handlers.UploadProjectImage)
			imageGroup.POST("/job-image", handlers.UploadJobImage)
			imageGroup.POST("/partners", handlers.UploadPartnerImage)
			imageGroup.POST("/carousel", handlers.UploadCarouselImage)
		}
	}

	analyticsAuthGroup := router.Group("/analytics")
	analyticsAuthGroup.Use(middlewares.AuthMiddleware())
	{
		analyticsAuthGroup.GET("", handlers.GetAnalytics)
		analyticsAuthGroup.GET("/count", handlers.CountAnalytics)
	}
	return
}

// Run defines the router endpoints and starts the server
func Run() {
	geoFilePath := fmt.Sprintf("./%s", geoKey)

	err := utils.DownloadFromS3(geoKey, geoFilePath)
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Failed to download geoip2 service file from s3")
	}

	db, err := geoip2.Open(geoFilePath)
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Unable to start geoip2 service")
	}
	defer db.Close()

	router := setupRouter(db)

	err = router.Run()
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Unable to start web server")
	}

	utils.GetLogger().Debug("Web server started ...")
}

// SetGeoFileKey gets the Geo File Key and stores it in memory
func SetGeoFileKey(key string) {
	geoKey = key
}
