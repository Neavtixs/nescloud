package middleware

import (
	"net/http"
	"strings"

	"nescloud/backend-app/internal/dto"
	"nescloud/backend-app/internal/errs"
	"nescloud/backend-app/internal/helper"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Authorization(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logEntry := helper.NewLog(log, c)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logEntry.Warn("authorization header missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			logEntry.Warn("invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
			return
		}

		userID, err := helper.ValidateAccessToken(tokenString)
		if err != nil {
			logEntry.WithField("layer", "middleware").Warn("invalid access token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
			return
		}

		logEntry.WithField("user_id", userID).Info("authorization success")
		c.Set("user_id", userID)
		c.Next()
	}
}
