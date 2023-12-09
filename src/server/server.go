package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"portfolio-cms-server/server/handlers"
	"portfolio-cms-server/server/middlewares"
	"portfolio-cms-server/utils"
)

func setupRouter() (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	router.Use(middlewares.Logger(utils.GetLogger()), gin.Recovery())
	router.Use(middlewares.CORS())

	router.GET("/healths", handlers.HealthCheck)
	router.GET("/metrics", handlers.Metrics)
	router.GET("/users/basic-info", handlers.GetBasicInfo)
	router.GET("/users/skills", handlers.GetSkills)
	router.GET("/users/jobs-and-projects", handlers.GetJobsAndProjects)
	router.GET("/users/socials", handlers.GetSocials)
	router.POST("/auth/login", handlers.Login)

	fileAuthGroup := router.Group("/files")
	fileAuthGroup.Use(middlewares.AuthMiddleware())
	{
		fileAuthGroup.POST("/cv", handlers.UploadCV)

		imageGroup := fileAuthGroup.Group("")
		imageGroup.Use(middlewares.ImageContentTypeMiddleware())
		{
			imageGroup.POST("/project-image", handlers.UploadProjectImage)
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
	router := setupRouter()

	err := router.Run()
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Unable to start web server")
	}

	utils.GetLogger().Debug("Web server started ...")
}
