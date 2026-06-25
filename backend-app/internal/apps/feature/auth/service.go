package auth

import (
	"database/sql"
	"fmt"
	"time"

	"nescloud/backend-app/internal/apps/domain/entity"
	"nescloud/backend-app/internal/apps/domain/repository"
	"nescloud/backend-app/internal/dto"
	"nescloud/backend-app/internal/helper"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
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

func (s *Service) Register(input *dto.InputAuthRegister) (*dto.ResultAuthRegister, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	userID := uuid.NewString()

	user := &entity.User{
		ID:        userID,
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashed),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.UserRepo.Insert(tx, input.Ctx, user); err != nil {
		return nil, err
	}

	quota := &entity.Quota{
		ID:          uuid.NewString(),
		UserID:      userID,
		MaxStorage:  1073741824,
		UsedStorage: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.QuotaRepo.Insert(tx, input.Ctx, quota); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	accessToken, err := helper.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.NewString()
	refreshTokenKey := fmt.Sprintf("refresh_token:%s", refreshToken)
	if err := s.Redis.Set(input.Ctx, refreshTokenKey, userID, 7*24*time.Hour).Err(); err != nil {
		return nil, err
	}

	return &dto.ResultAuthRegister{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
