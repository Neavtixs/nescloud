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

type FolderRepo struct{}

func NewFolderRepo() *FolderRepo {
	return &FolderRepo{}
}

func (r *FolderRepo) Insert(tx *sql.Tx, ctx context.Context, folder *entity.Folder) error {
	query := `
		INSERT INTO folders (id, owner_id, parent_folder_id, name, deleted_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.ExecContext(ctx, query,
		folder.ID, folder.OwnerID, folder.ParentFolderID, folder.Name,
		folder.DeletedAt, folder.CreatedAt, folder.UpdatedAt,
	)
	return err
}

func (r *FolderRepo) FindByID(tx *sql.Tx, ctx context.Context, id string) (*entity.Folder, error) {
	query := `
		SELECT id, owner_id, parent_folder_id, name, deleted_at, created_at, updated_at
		FROM folders
		WHERE id = $1
	`
	folder := &entity.Folder{}
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&folder.ID, &folder.OwnerID, &folder.ParentFolderID, &folder.Name,
		&folder.DeletedAt, &folder.CreatedAt, &folder.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrDataNotFound
		}
		return nil, err
	}
	return folder, nil
}

func (r *FolderRepo) FindByOwnerID(tx *sql.Tx, ctx context.Context, ownerID string, parentFolderID *string, limit, offset int) ([]entity.Folder, int, error) {
	var total int

	countQuery := `
		SELECT COUNT(*)
		FROM folders
		WHERE owner_id = $1 AND deleted_at IS NULL AND parent_folder_id IS NOT DISTINCT FROM $2
	`
	if err := tx.QueryRowContext(ctx, countQuery, ownerID, parentFolderID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, owner_id, parent_folder_id, name, deleted_at, created_at, updated_at
		FROM folders
		WHERE owner_id = $1 AND deleted_at IS NULL AND parent_folder_id IS NOT DISTINCT FROM $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := tx.QueryContext(ctx, query, ownerID, parentFolderID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var folders []entity.Folder
	for rows.Next() {
		var f entity.Folder
		if err := rows.Scan(&f.ID, &f.OwnerID, &f.ParentFolderID, &f.Name, &f.DeletedAt, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, 0, err
		}
		folders = append(folders, f)
	}
	if folders == nil {
		folders = []entity.Folder{}
	}
	return folders, total, nil
}

func (r *FolderRepo) UpdateName(tx *sql.Tx, ctx context.Context, id, name string) error {
	query := `
		UPDATE folders
		SET name = $1, updated_at = $2
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

func (r *FolderRepo) SoftDelete(tx *sql.Tx, ctx context.Context, id string, deletedAt time.Time) error {
	query := `
		UPDATE folders
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

func (r *FolderRepo) Restore(tx *sql.Tx, ctx context.Context, id string) error {
	query := `
		UPDATE folders
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

func (r *FolderRepo) HardDelete(tx *sql.Tx, ctx context.Context, id string) error {
	query := `DELETE FROM folders WHERE id = $1`
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

func (r *FolderRepo) FindDeletedByOwnerID(tx *sql.Tx, ctx context.Context, ownerID string) ([]entity.Folder, error) {
	query := `
		SELECT id, owner_id, parent_folder_id, name, deleted_at, created_at, updated_at
		FROM folders
		WHERE owner_id = $1 AND deleted_at IS NOT NULL
		ORDER BY deleted_at DESC
	`
	rows, err := tx.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []entity.Folder
	for rows.Next() {
		var f entity.Folder
		if err := rows.Scan(&f.ID, &f.OwnerID, &f.ParentFolderID, &f.Name, &f.DeletedAt, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}
	if folders == nil {
		folders = []entity.Folder{}
	}
	return folders, nil
}

func (r *FolderRepo) CountChildren(tx *sql.Tx, ctx context.Context, parentID string) (int, error) {
	query := `SELECT COUNT(*) FROM folders WHERE parent_folder_id = $1 AND deleted_at IS NULL`
	var count int
	err := tx.QueryRowContext(ctx, query, parentID).Scan(&count)
	return count, err
}

func (r *FolderRepo) FindSubfolderIDsRecursive(tx *sql.Tx, ctx context.Context, folderID string) ([]string, error) {
	query := `
		WITH RECURSIVE subfolders AS (
			SELECT id FROM folders WHERE id = $1 AND deleted_at IS NULL
			UNION ALL
			SELECT f.id FROM folders f
			INNER JOIN subfolders s ON f.parent_folder_id = s.id
			WHERE f.deleted_at IS NULL
		)
		SELECT id FROM subfolders
	`
	rows, err := tx.QueryContext(ctx, query, folderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if ids == nil {
		ids = []string{}
	}
	return ids, nil
}

func (r *FolderRepo) BulkHardDelete(tx *sql.Tx, ctx context.Context, ids []string) error {
	query := `DELETE FROM folders WHERE id = ANY($1)`
	_, err := tx.ExecContext(ctx, query, pq.Array(ids))
	return err
}
