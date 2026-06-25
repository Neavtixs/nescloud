package file

import (
	"database/sql"

	"nescloud/backend-app/internal/apps/domain/repository"
	"nescloud/backend-app/internal/apps/storage"
)

type Service struct {
	DB        *sql.DB
	FileRepo  *repository.FileRepo
	QuotaRepo *repository.QuotaRepo
	Storage   *storage.Storage
}

func NewService(db *sql.DB, fileRepo *repository.FileRepo, quotaRepo *repository.QuotaRepo, store *storage.Storage) *Service {
	return &Service{
		DB:        db,
		FileRepo:  fileRepo,
		QuotaRepo: quotaRepo,
		Storage:   store,
	}
}
