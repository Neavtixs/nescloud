package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"nescloud/backend-app/internal/apps/domain/entity"
	"nescloud/backend-app/internal/apps/domain/repository"
	"nescloud/backend-app/internal/dto"
	"nescloud/backend-app/internal/errs"
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

func (s *Service) setRefreshToken(ctx context.Context, userID string) (string, error) {
	refreshToken := uuid.NewString()
	refreshTokenKey := fmt.Sprintf("refresh_token:%s", refreshToken)
	/*
		refreshTTL := 7 * 24 * time.Hour
		if d, err := strconv.Atoi(os.Getenv("REFRESH_EXP_DAYS")); err == nil && d > 0 {
			refreshTTL = time.Duration(d) * 24 * time.Hour
		}
	*/
	if err := s.Redis.Set(ctx, refreshTokenKey, userID, 90*time.Second).Err(); err != nil {
		return "", err
	}
	return refreshToken, nil
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

	refreshToken, err := s.setRefreshToken(input.Ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.ResultAuthRegister{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Login(input *dto.InputAuthLogin) (*dto.ResultAuthLogin, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := s.UserRepo.FindByEmail(tx, input.Ctx, input.Email)
	if err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			return nil, errs.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	accessToken, err := helper.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.setRefreshToken(input.Ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.ResultAuthLogin{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Logout(input *dto.InputAuthLogout) error {
	if input.RefreshToken == "" {
		return nil
	}

	refreshTokenKey := fmt.Sprintf("refresh_token:%s", input.RefreshToken)
	_ = s.Redis.Del(input.Ctx, refreshTokenKey).Err()

	return nil
}

func (s *Service) Refresh(input *dto.InputAuthRefresh) (*dto.ResultAuthRefresh, error) {
	refreshTokenKey := fmt.Sprintf("refresh_token:%s", input.RefreshToken)

	userID, err := s.Redis.Get(input.Ctx, refreshTokenKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errs.ErrInvalidAccessToken
		}
		return nil, err
	}

	accessToken, err := helper.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	s.Redis.Del(input.Ctx, refreshTokenKey)

	newRefreshToken, err := s.setRefreshToken(input.Ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.ResultAuthRefresh{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) Me(input *dto.InputAuthMe) (*dto.ResultAuthMe, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := s.UserRepo.FindByID(tx, input.Ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.ResultAuthMe{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
