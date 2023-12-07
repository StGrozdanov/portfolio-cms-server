package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	handlers2 "portfolio-cms-server/server/handlers"
	"portfolio-cms-server/server/interceptors"
	"portfolio-cms-server/server/middlewares"
	"portfolio-cms-server/utils"
)

func setupRouter() (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	router.Use(middlewares.Logger(utils.GetLogger()), gin.Recovery())
	router.Use(interceptors.Interceptor())
	return
}

// Run defines the router endpoints and starts the server
func Run() {
	router := setupRouter()
	router.GET("/healths", handlers2.HealthCheck)
	router.GET("/metrics", handlers2.Metrics)

	err := router.Run()
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Unable to start web server")
	}
	utils.GetLogger().Debug("Web server started ...")
}
