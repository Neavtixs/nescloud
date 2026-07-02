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

type FileRepo struct{}

func NewFileRepo() *FileRepo {
	return &FileRepo{}
}

func (r *FileRepo) Insert(tx *sql.Tx, ctx context.Context, file *entity.File) error {
	query := `
		INSERT INTO files (id, owner_id, folder_id, original_name, storage_key, mime_type, extension, size, checksum, upload_status, deleted_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := tx.ExecContext(ctx, query,
		file.ID, file.OwnerID, file.FolderID, file.OriginalName, file.StorageKey,
		file.MimeType, file.Extension, file.Size, file.Checksum, file.UploadStatus,
		file.DeletedAt, file.CreatedAt, file.UpdatedAt,
	)
	return err
}

func (r *FileRepo) FindByID(tx *sql.Tx, ctx context.Context, id string) (*entity.File, error) {
	query := `
		SELECT id, owner_id, folder_id, original_name, storage_key, mime_type, extension, size, checksum, upload_status, deleted_at, created_at, updated_at
		FROM files
		WHERE id = $1
	`
	file := &entity.File{}
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&file.ID, &file.OwnerID, &file.FolderID, &file.OriginalName, &file.StorageKey,
		&file.MimeType, &file.Extension, &file.Size, &file.Checksum, &file.UploadStatus,
		&file.DeletedAt, &file.CreatedAt, &file.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrDataNotFound
		}
		return nil, err
	}
	return file, nil
}

func (r *FileRepo) FindByOwnerIDAndFolderID(tx *sql.Tx, ctx context.Context, ownerID string, folderID *string, search string, limit, offset int) ([]entity.File, int, error) {
	var total int
	var rows *sql.Rows
	var err error

	if search != "" {
		countQuery := `
			SELECT COUNT(*)
			FROM files
			WHERE owner_id = $1 AND folder_id IS NOT DISTINCT FROM $2 AND deleted_at IS NULL AND original_name ILIKE '%' || $3 || '%'
		`
		if err = tx.QueryRowContext(ctx, countQuery, ownerID, folderID, search).Scan(&total); err != nil {
			return nil, 0, err
		}

		query := `
			SELECT id, owner_id, folder_id, original_name, storage_key, mime_type, extension, size, checksum, upload_status, deleted_at, created_at, updated_at
			FROM files
			WHERE owner_id = $1 AND folder_id IS NOT DISTINCT FROM $2 AND deleted_at IS NULL AND original_name ILIKE '%' || $3 || '%'
			ORDER BY created_at DESC
			LIMIT $4 OFFSET $5
		`
		rows, err = tx.QueryContext(ctx, query, ownerID, folderID, search, limit, offset)
	} else {
		countQuery := `
			SELECT COUNT(*)
			FROM files
			WHERE owner_id = $1 AND folder_id IS NOT DISTINCT FROM $2 AND deleted_at IS NULL
		`
		if err = tx.QueryRowContext(ctx, countQuery, ownerID, folderID).Scan(&total); err != nil {
			return nil, 0, err
		}

		query := `
			SELECT id, owner_id, folder_id, original_name, storage_key, mime_type, extension, size, checksum, upload_status, deleted_at, created_at, updated_at
			FROM files
			WHERE owner_id = $1 AND folder_id IS NOT DISTINCT FROM $2 AND deleted_at IS NULL
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`
		rows, err = tx.QueryContext(ctx, query, ownerID, folderID, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var files []entity.File
	for rows.Next() {
		var f entity.File
		if err := rows.Scan(&f.ID, &f.OwnerID, &f.FolderID, &f.OriginalName, &f.StorageKey,
			&f.MimeType, &f.Extension, &f.Size, &f.Checksum, &f.UploadStatus,
			&f.DeletedAt, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		files = append(files, f)
	}
	if files == nil {
		files = []entity.File{}
	}
	return files, total, nil
}

func (r *FileRepo) UpdateName(tx *sql.Tx, ctx context.Context, id, name string) error {
	query := `
		UPDATE files
		SET original_name = $1, updated_at = $2
		WHERE id = $3
	`
	result, err := tx.ExecContext(ctx, query, name, time.Now(), id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errs.ErrDataNotFound
	}
	return nil
}

func (r *FileRepo) UpdateStatus(tx *sql.Tx, ctx context.Context, id, status string) error {
	query := `
		UPDATE files
		SET upload_status = $1, updated_at = $2
		WHERE id = $3
	`
	result, err := tx.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errs.ErrDataNotFound
	}
	return nil
}

func (r *FileRepo) SoftDelete(tx *sql.Tx, ctx context.Context, id string, deletedAt time.Time) error {
	query := `
		UPDATE files
		SET deleted_at = $1, updated_at = $2
		WHERE id = $3
	`
	result, err := tx.ExecContext(ctx, query, deletedAt, time.Now(), id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errs.ErrDataNotFound
	}
	return nil
}

func (r *FileRepo) Restore(tx *sql.Tx, ctx context.Context, id string) error {
	query := `
		UPDATE files
		SET deleted_at = NULL, updated_at = $1
		WHERE id = $2 AND deleted_at IS NOT NULL
	`
	result, err := tx.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errs.ErrDataNotFound
	}
	return nil
}

func (r *FileRepo) HardDelete(tx *sql.Tx, ctx context.Context, id string) error {
	query := `DELETE FROM files WHERE id = $1`
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errs.ErrDataNotFound
	}
	return nil
}

func (r *FileRepo) FindDeletedByOwnerID(tx *sql.Tx, ctx context.Context, ownerID string) ([]entity.File, error) {
	query := `
		SELECT id, owner_id, folder_id, original_name, storage_key, mime_type, extension, size, checksum, upload_status, deleted_at, created_at, updated_at
		FROM files
		WHERE owner_id = $1 AND deleted_at IS NOT NULL
		ORDER BY deleted_at DESC
	`
	rows, err := tx.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []entity.File
	for rows.Next() {
		var f entity.File
		if err := rows.Scan(&f.ID, &f.OwnerID, &f.FolderID, &f.OriginalName, &f.StorageKey,
			&f.MimeType, &f.Extension, &f.Size, &f.Checksum, &f.UploadStatus,
			&f.DeletedAt, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	if files == nil {
		files = []entity.File{}
	}
	return files, nil
}

func (r *FileRepo) FindFilesByFolderIDs(tx *sql.Tx, ctx context.Context, folderIDs []string) ([]entity.File, error) {
	query := `
		SELECT id, owner_id, folder_id, original_name, storage_key, mime_type, extension, size, checksum, upload_status, deleted_at, created_at, updated_at
		FROM files
		WHERE folder_id = ANY($1) AND deleted_at IS NULL
	`
	rows, err := tx.QueryContext(ctx, query, pq.Array(folderIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []entity.File
	for rows.Next() {
		var f entity.File
		if err := rows.Scan(&f.ID, &f.OwnerID, &f.FolderID, &f.OriginalName, &f.StorageKey,
			&f.MimeType, &f.Extension, &f.Size, &f.Checksum, &f.UploadStatus,
			&f.DeletedAt, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	if files == nil {
		files = []entity.File{}
	}
	return files, nil
}

func (r *FileRepo) BulkSoftDelete(tx *sql.Tx, ctx context.Context, ids []string, deletedAt time.Time) error {
	query := `
		UPDATE files
		SET deleted_at = $1, updated_at = $2
		WHERE id = ANY($3)
	`
	_, err := tx.ExecContext(ctx, query, deletedAt, time.Now(), pq.Array(ids))
	return err
}

func (r *FileRepo) BulkHardDelete(tx *sql.Tx, ctx context.Context, ids []string) error {
	query := `DELETE FROM files WHERE id = ANY($1)`
	_, err := tx.ExecContext(ctx, query, pq.Array(ids))
	return err
}
