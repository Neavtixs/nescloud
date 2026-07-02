package repository

import (
	"context"
	"database/sql"
	"errors"

	"nescloud/backend-app/internal/apps/domain/entity"
	"nescloud/backend-app/internal/errs"
)

type PublicLinkRepo struct{}

func NewPublicLinkRepo() *PublicLinkRepo {
	return &PublicLinkRepo{}
}

func (r *PublicLinkRepo) Insert(tx *sql.Tx, ctx context.Context, link *entity.PublicLink) error {
	query := `
		INSERT INTO public_links (id, file_id, token, expired_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := tx.ExecContext(ctx, query, link.ID, link.FileID, link.Token, link.ExpiredAt, link.CreatedAt)
	return err
}

func (r *PublicLinkRepo) FindByToken(tx *sql.Tx, ctx context.Context, token string) (*entity.PublicLink, error) {
	query := `
		SELECT id, file_id, token, expired_at, created_at
		FROM public_links
		WHERE token = $1
	`
	link := &entity.PublicLink{}
	err := tx.QueryRowContext(ctx, query, token).Scan(
		&link.ID, &link.FileID, &link.Token, &link.ExpiredAt, &link.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrDataNotFound
		}
		return nil, err
	}
	return link, nil
}

func (r *PublicLinkRepo) FindByFileID(tx *sql.Tx, ctx context.Context, fileID string) (*entity.PublicLink, error) {
	query := `
		SELECT id, file_id, token, expired_at, created_at
		FROM public_links
		WHERE file_id = $1
	`
	link := &entity.PublicLink{}
	err := tx.QueryRowContext(ctx, query, fileID).Scan(
		&link.ID, &link.FileID, &link.Token, &link.ExpiredAt, &link.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrDataNotFound
		}
		return nil, err
	}
	return link, nil
}

func (r *PublicLinkRepo) Delete(tx *sql.Tx, ctx context.Context, fileID string) error {
	query := `DELETE FROM public_links WHERE file_id = $1`
	result, err := tx.ExecContext(ctx, query, fileID)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errs.ErrDataNotFound
	}
	return nil
}

func (r *PublicLinkRepo) FindAllByOwnerID(tx *sql.Tx, ctx context.Context, ownerID string) ([]entity.PublicLinkFile, error) {
	query := `
		SELECT pl.id, pl.file_id, pl.token, pl.expired_at, pl.created_at,
		       f.original_name, f.mime_type, f.size
		FROM public_links pl
		INNER JOIN files f ON f.id = pl.file_id
		WHERE f.owner_id = $1 AND f.deleted_at IS NULL
		ORDER BY pl.created_at DESC
	`
	rows, err := tx.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.PublicLinkFile
	for rows.Next() {
		var item entity.PublicLinkFile
		if err := rows.Scan(
			&item.ID, &item.FileID, &item.Token, &item.ExpiredAt, &item.CreatedAt,
			&item.OriginalName, &item.MimeType, &item.Size,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	if result == nil {
		result = []entity.PublicLinkFile{}
	}
	return result, nil
}
