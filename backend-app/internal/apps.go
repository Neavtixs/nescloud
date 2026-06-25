package internal

import (
	"database/sql"

	"nescloud/backend-app/configs"
	"nescloud/backend-app/internal/apps/domain/repository"
	"nescloud/backend-app/internal/apps/feature/auth"
	"nescloud/backend-app/internal/apps/feature/file"
	"nescloud/backend-app/internal/apps/feature/folder"
	"nescloud/backend-app/internal/apps/feature/publiclink"
	"nescloud/backend-app/internal/apps/feature/trash"
	"nescloud/backend-app/internal/apps/storage"
	"nescloud/backend-app/internal/route"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func InitDependencies(db *sql.DB, rdb *redis.Client, validate *validator.Validate, log *logrus.Logger) *route.Handler {
	userRepo := repository.NewUserRepo()
	quotaRepo := repository.NewQuotaRepo()
	folderRepo := repository.NewFolderRepo()
	fileRepo := repository.NewFileRepo()
	publicLinkRepo := repository.NewPublicLinkRepo()
	auditLogRepo := repository.NewAuditLogRepo()

	s3cfg := configs.NewS3()
	store := storage.NewStorage(s3cfg)

	authService := auth.NewService(db, rdb, userRepo, quotaRepo)
	authHandler := auth.NewHandler(authService, validate, log)

	folderService := folder.NewService(db, folderRepo)
	_ = folder.NewHandler(folderService, validate, log)

	fileService := file.NewService(db, fileRepo, quotaRepo, store)
	_ = file.NewHandler(fileService, validate, log)

	publicLinkService := publiclink.NewService(db, publicLinkRepo)
	_ = publiclink.NewHandler(publicLinkService, validate, log)

	trashService := trash.NewService(db, folderRepo, fileRepo, quotaRepo, store)
	_ = trash.NewHandler(trashService, validate, log)

	_ = auditLogRepo

	return route.NewHandler(authHandler)
}
