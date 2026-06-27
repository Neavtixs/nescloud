package route

import (
	"nescloud/backend-app/internal/apps/feature/auth"
	"nescloud/backend-app/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Auth *auth.Handler
	Log  *logrus.Logger
}

func NewHandler(auth *auth.Handler, log *logrus.Logger) *Handler {
	return &Handler{Auth: auth, Log: log}
}

func (h *Handler) SetupRoute(app *gin.Engine) {
	api := app.Group("/api")

	public := api.Group("")
	{
		public.POST("/auth/register", h.Auth.RegisterHandler)
		public.POST("/auth/login", h.Auth.LoginHandler)
	}

	user := api.Group("", middleware.Authorization(h.Log))
	{
		user.POST("/auth/logout", h.Auth.LogoutHandler)
		user.GET("/auth/me", h.Auth.MeHandler)
	}
}
