package publiclink

import (
	"database/sql"

	"nescloud/backend-app/internal/apps/domain/repository"
)

type Service struct {
	DB             *sql.DB
	PublicLinkRepo *repository.PublicLinkRepo
}

func NewService(db *sql.DB, publicLinkRepo *repository.PublicLinkRepo) *Service {
	return &Service{
		DB:             db,
		PublicLinkRepo: publicLinkRepo,
	}
}
