package auth

import (
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
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	DB        *sql.DB
	Redis     *redis.Client
	UserRepo  *repository.UserRepo
	QuotaRepo *repository.QuotaRepo
	Log       *logrus.Logger
}

func NewService(db *sql.DB, rdb *redis.Client, userRepo *repository.UserRepo, quotaRepo *repository.QuotaRepo, log *logrus.Logger) *Service {
	return &Service{
		DB:        db,
		Redis:     rdb,
		UserRepo:  userRepo,
		QuotaRepo: quotaRepo,
		Log:       log,
	}
}

func (s *Service) Register(input *dto.InputAuthRegister) (*dto.ResultAuthRegister, error) {
	s.Log.WithField("email", input.Email).Info("register initiated")

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
	s.Log.WithField("user_id", userID).Info("user inserted into db")

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
	s.Log.WithField("user_id", userID).Info("quota created for user")

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	s.Log.WithField("user_id", userID).Info("transaction committed")

	accessToken, err := helper.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.NewString()
	refreshTokenKey := fmt.Sprintf("refresh_token:%s", refreshToken)
	if err := s.Redis.Set(input.Ctx, refreshTokenKey, userID, 7*24*time.Hour).Err(); err != nil {
		return nil, err
	}
	s.Log.WithField("user_id", userID).Info("refresh token saved to redis")

	s.Log.WithField("user_id", userID).Info("register completed")

	return &dto.ResultAuthRegister{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Login(input *dto.InputAuthLogin) (*dto.ResultAuthLogin, error) {
	s.Log.WithField("email", input.Email).Info("login initiated")

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
	s.Log.WithField("user_id", user.ID).Info("user found")

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errs.ErrInvalidCredentials
	}
	s.Log.WithField("user_id", user.ID).Info("password verified")

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	s.Log.WithField("user_id", user.ID).Info("transaction committed")

	accessToken, err := helper.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.NewString()
	refreshTokenKey := fmt.Sprintf("refresh_token:%s", refreshToken)
	if err := s.Redis.Set(input.Ctx, refreshTokenKey, user.ID, 7*24*time.Hour).Err(); err != nil {
		return nil, err
	}
	s.Log.WithField("user_id", user.ID).Info("refresh token saved to redis")

	s.Log.WithField("user_id", user.ID).Info("login completed")

	return &dto.ResultAuthLogin{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Logout(input *dto.InputAuthLogout) error {
	s.Log.Info("logout initiated")

	if input.RefreshToken == "" {
		s.Log.Warn("no refresh token provided for logout")
		return nil
	}

	refreshTokenKey := fmt.Sprintf("refresh_token:%s", input.RefreshToken)
	if err := s.Redis.Del(input.Ctx, refreshTokenKey).Err(); err != nil {
		s.Log.WithError(err).Warn("failed to delete refresh token from redis")
	} else {
		s.Log.Info("refresh token deleted from redis")
	}

	return nil
}

func (s *Service) Me(input *dto.InputAuthMe) (*dto.ResultAuthMe, error) {
	s.Log.WithField("user_id", input.UserID).Info("me initiated")

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

	s.Log.WithField("user_id", input.UserID).Info("me completed")

	return &dto.ResultAuthMe{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
