package route

import (
	"os"
	"strings"

	"nescloud/backend-app/internal/apps/feature/auth"
	"nescloud/backend-app/internal/apps/feature/folder"
	"nescloud/backend-app/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Auth   *auth.Handler
	Folder *folder.Handler
	Log    *logrus.Logger
}

func NewHandler(auth *auth.Handler, folder *folder.Handler, log *logrus.Logger) *Handler {
	return &Handler{Auth: auth, Folder: folder, Log: log}
}

func (h *Handler) SetupRoute(app *gin.Engine) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(os.Getenv("CORS_ALLOW_ORIGIN"), ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	api := app.Group("/api")

	public := api.Group("")
	{
		public.POST("/auth/register", h.Auth.RegisterHandler)
		public.POST("/auth/login", h.Auth.LoginHandler)
		public.POST("/auth/refresh", h.Auth.RefreshHandler)
	}

	user := api.Group("", middleware.Authorization(h.Log))
	{
		user.POST("/auth/logout", h.Auth.LogoutHandler)
		user.GET("/auth/me", h.Auth.MeHandler)

		folders := user.Group("/folders")
		{
			folders.POST("", h.Folder.CreateFolderHandler)
			folders.GET("", h.Folder.ListFoldersHandler)
			folders.PATCH("/:id", h.Folder.UpdateFolderHandler)
			folders.DELETE("/:id", h.Folder.DeleteFolderHandler)
		}

		trash := user.Group("/trash")
		{
			trash.GET("/folders", h.Folder.ListTrashHandler)
			trash.POST("/folders/:id/restore", h.Folder.RestoreFolderHandler)
			trash.DELETE("/folders/:id", h.Folder.PermanentDeleteFolderHandler)
		}
	}
}
