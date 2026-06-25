package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func NewLog(log *logrus.Logger, c *gin.Context) *logrus.Entry {
	return log.WithFields(logrus.Fields{
		"request_id": c.GetString("request_id"),
		"path":       c.Request.URL.Path,
		"method":     c.Request.Method,
	})
}

func ValidationMsg(err error) map[string]string {
	res := make(map[string]string)

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			res[fe.Field()] = fe.Tag()
		}
	}

	return res
}
