package repository

import (
	"context"
	"database/sql"
	"errors"

	"nescloud/backend-app/internal/apps/domain/entity"
	"nescloud/backend-app/internal/errs"
)

type QuotaRepo struct{}

func NewQuotaRepo() *QuotaRepo {
	return &QuotaRepo{}
}

func (r *QuotaRepo) Insert(tx *sql.Tx, ctx context.Context, quota *entity.Quota) error {
	query := `
		INSERT INTO quotas (id, user_id, max_storage, used_storage, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.ExecContext(ctx, query,
		quota.ID, quota.UserID, quota.MaxStorage, quota.UsedStorage,
		quota.CreatedAt, quota.UpdatedAt,
	)
	return err
}

func (r *QuotaRepo) FindByUserID(tx *sql.Tx, ctx context.Context, userID string) (*entity.Quota, error) {
	query := `
		SELECT id, user_id, max_storage, used_storage, created_at, updated_at
		FROM quotas
		WHERE user_id = $1
	`
	quota := &entity.Quota{}
	err := tx.QueryRowContext(ctx, query, userID).Scan(
		&quota.ID, &quota.UserID, &quota.MaxStorage, &quota.UsedStorage,
		&quota.CreatedAt, &quota.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrDataNotFound
		}
		return nil, err
	}
	return quota, nil
}
