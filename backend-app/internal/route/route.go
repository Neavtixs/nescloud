package route

import (
	"github.com/gin-gonic/gin"
)

func SetupRoute(app *gin.Engine) {
	api := app.Group("/api")
	_ = api
}
