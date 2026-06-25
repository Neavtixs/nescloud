package folder

import (
	"database/sql"

	"nescloud/backend-app/internal/apps/domain/repository"
)

type Service struct {
	DB         *sql.DB
	FolderRepo *repository.FolderRepo
}

func NewService(db *sql.DB, folderRepo *repository.FolderRepo) *Service {
	return &Service{
		DB:         db,
		FolderRepo: folderRepo,
	}
}
