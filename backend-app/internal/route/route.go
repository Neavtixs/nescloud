package route

import (
	"nescloud/backend-app/internal/apps/feature/auth"
	"nescloud/backend-app/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Auth *auth.Handler
}

func NewHandler(auth *auth.Handler) *Handler {
	return &Handler{Auth: auth}
}

func (h *Handler) SetupRoute(app *gin.Engine) {
	api := app.Group("/api")

	public := api.Group("")
	{
		public.POST("/auth/register", h.Auth.RegisterHandler)
		public.POST("/auth/login", h.Auth.LoginHandler)
	}

	user := api.Group("", middleware.Authorization())
	{
		user.POST("/auth/logout", h.Auth.LogoutHandler)
	}
}
