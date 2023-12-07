package handlers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"portfolio-cms-server/internal"
	"portfolio-cms-server/utils"
)

func GetBasicInfo(ginCtx *gin.Context) {
	basicUserInfo, err := internal.GetBasicUserInfo()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, basicUserInfo)
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting basic user info from the database")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, basicUserInfo)
}
