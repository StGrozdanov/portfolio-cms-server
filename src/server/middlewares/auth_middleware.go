package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"portfolio-cms-server/utils"
)

// AuthMiddleware checks for valid x-authorization access token in the request headers
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header["X-Authorization"][0]

		if len(token) == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, isValid, err := utils.ParseJWT(token)
		if err != nil || !isValid {
			ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]interface{}{"message": "Invalid token"})
			return
		}

		if claims.Role != "administrator" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]interface{}{"message": "Invalid permissions"})
			return
		}

		ctx.Next()
	}
}
