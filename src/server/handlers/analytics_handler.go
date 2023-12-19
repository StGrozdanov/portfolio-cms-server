package handlers

import (
	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
	"net/http"
	"portfolio-cms-server/internal/analytics"
	"portfolio-cms-server/utils"
	"strings"
)

func GetAnalytics(ginCtx *gin.Context) {
	params := ginCtx.Request.URL.RawQuery

	query := strings.Split(params, "=")
	if len(query) == 1 {
		ginCtx.AddParam(query[0], "")
	} else if len(query) == 2 {
		ginCtx.AddParam(query[0], query[1])
	} else {
		ginCtx.JSON(
			http.StatusBadRequest,
			map[string]interface{}{"message": "You should provide exactly 1 query param to this endpoint"},
		)
		return
	}

	analyticsResults, err := analytics.Get(ginCtx.Params[0])
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		} else if strings.Contains(err.Error(), "param") {
			ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"message": err.Error()})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting analytics from the database")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, analyticsResults)
}

func CountAnalytics(ginCtx *gin.Context) {
	count, err := analytics.Count()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			ginCtx.JSON(http.StatusOK, map[string]interface{}{})
			return
		}

		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on getting analytics count from the database")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, map[string]interface{}{"count": count})
}

func Track(ginCtx *gin.Context, db *geoip2.Reader) {
	requestBody := analytics.UserDevice{}

	if err := ginCtx.ShouldBind(&requestBody); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	if _, err := validator.ValidateStruct(requestBody); err != nil {
		ginCtx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid parameters"})
		return
	}

	err := analytics.Track(db, ginCtx, requestBody.DeviceType)
	if err != nil {
		utils.
			GetLogger().
			WithFields(log.Fields{"error": err.Error()}).
			Error("Error on tracking user request")

		ginCtx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ginCtx.JSON(http.StatusOK, map[string]interface{}{"status": "recorded"})
}
