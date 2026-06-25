package middleware

import (
	"net/http"
	"os"
	"strings"

	"nescloud/backend-app/internal/dto"
	"nescloud/backend-app/internal/errs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if strings.TrimSpace(secret) == "" {
			secret = "default-secret"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errs.ErrInvalidAccessToken
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
