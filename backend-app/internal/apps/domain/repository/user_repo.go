package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"nescloud/backend-app/internal/apps/domain/entity"
	"nescloud/backend-app/internal/errs"

	"github.com/lib/pq"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Insert(tx *sql.Tx, ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.ExecContext(ctx, query,
		user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return errs.ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UserRepo) FindByEmail(tx *sql.Tx, ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &entity.User{}
	err := tx.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrDataNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) FindByID(tx *sql.Tx, ctx context.Context, id string) (*entity.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &entity.User{}
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrDataNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) Update(tx *sql.Tx, ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := tx.ExecContext(ctx, query,
		user.Name, user.Email, user.Password, time.Now(), user.ID,
	)
	return err
}
