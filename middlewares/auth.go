package middlewares

import (
	"dev_community_server/crypto"
	"dev_community_server/utils"
	"github.com/gin-gonic/gin"
	log "github.com/shyunku-libraries/go-logger"
	"net/http"
)

func TokenMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		accessToken, err := utils.GetBearerToken(authHeader)
		if err != nil {
			log.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		atClaims, atErr := crypto.ValidateAccessToken(*accessToken)

		if atErr == nil {
			c.Set("uuid", atClaims.Uuid)
			log.Info("valid token / [uuid]:", atClaims.Uuid)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
