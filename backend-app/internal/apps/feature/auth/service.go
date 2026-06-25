package auth

import (
	"database/sql"

	"nescloud/backend-app/internal/apps/domain/repository"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	DB        *sql.DB
	Redis     *redis.Client
	UserRepo  *repository.UserRepo
	QuotaRepo *repository.QuotaRepo
}

func NewService(db *sql.DB, rdb *redis.Client, userRepo *repository.UserRepo, quotaRepo *repository.QuotaRepo) *Service {
	return &Service{
		DB:        db,
		Redis:     rdb,
		UserRepo:  userRepo,
		QuotaRepo: quotaRepo,
	}
}
