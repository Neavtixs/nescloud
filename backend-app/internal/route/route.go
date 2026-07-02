package route

import (
	"os"
	"strings"

	"nescloud/backend-app/internal/apps/feature/auth"
	"nescloud/backend-app/internal/apps/feature/file"
	"nescloud/backend-app/internal/apps/feature/folder"
	"nescloud/backend-app/internal/apps/feature/publiclink"
	"nescloud/backend-app/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Auth       *auth.Handler
	Folder     *folder.Handler
	File       *file.Handler
	PublicLink *publiclink.Handler
	Log        *logrus.Logger
}

func NewHandler(auth *auth.Handler, folder *folder.Handler, file *file.Handler, publicLink *publiclink.Handler, log *logrus.Logger) *Handler {
	return &Handler{Auth: auth, Folder: folder, File: file, PublicLink: publicLink, Log: log}
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
		public.GET("/public/:token", h.PublicLink.AccessHandler)
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
			trash.DELETE("", h.Folder.EmptyTrashHandler)
			trash.GET("/folders", h.Folder.ListTrashHandler)
			trash.POST("/folders/:id/restore", h.Folder.RestoreFolderHandler)
			trash.DELETE("/folders/:id", h.Folder.PermanentDeleteFolderHandler)

			trash.GET("/files", h.File.ListTrashFilesHandler)
			trash.POST("/files/:id/restore", h.File.RestoreFileHandler)
			trash.DELETE("/files/:id", h.File.PermanentDeleteFileHandler)
		}

		files := user.Group("/files")
		{
			files.POST("/init-upload", h.File.InitUploadHandler)
			files.POST("/complete", h.File.CompleteUploadHandler)
			files.GET("", h.File.ListFilesHandler)
			files.GET("/public-links", h.PublicLink.ListLinksHandler)
			files.GET("/:id/download", h.File.DownloadHandler)
			files.GET("/:id", h.File.FileDetailHandler)
			files.PATCH("/:id", h.File.RenameHandler)
			files.DELETE("/:id", h.File.DeleteFileHandler)
			files.POST("/:id/public-link", h.PublicLink.GenerateHandler)
			files.DELETE("/:id/public-link", h.PublicLink.RevokeHandler)
		}
	}
}
