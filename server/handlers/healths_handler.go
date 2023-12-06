package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"portfolio-cms-server/database"
)

type healthResponse struct {
	AppStatus      string `json:"AppStatus"`
	DatabaseStatus string `json:"Database"`
}

var health healthResponse

func HealthCheck(ginCtx *gin.Context) {
	checkDB(&health)
	ginCtx.JSON(http.StatusOK, health)
}

func checkDB(response *healthResponse) {
	err := database.Ping()
	if err != nil {
		response.AppStatus = "Unhealthy"
		response.DatabaseStatus = err.Error()
	} else {
		response.AppStatus = "Healthy"
		response.DatabaseStatus = "Healthy"
	}
	database.CloseConnection()
}
