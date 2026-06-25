package trash

import (
	"database/sql"

	"nescloud/backend-app/internal/apps/domain/repository"
	"nescloud/backend-app/internal/apps/storage"
)

type Service struct {
	DB         *sql.DB
	FolderRepo *repository.FolderRepo
	FileRepo   *repository.FileRepo
	QuotaRepo  *repository.QuotaRepo
	Storage    *storage.Storage
}

func NewService(db *sql.DB, folderRepo *repository.FolderRepo, fileRepo *repository.FileRepo, quotaRepo *repository.QuotaRepo, store *storage.Storage) *Service {
	return &Service{
		DB:         db,
		FolderRepo: folderRepo,
		FileRepo:   fileRepo,
		QuotaRepo:  quotaRepo,
		Storage:    store,
	}
}
