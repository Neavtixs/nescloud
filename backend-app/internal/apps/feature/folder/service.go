package folder

import (
	"database/sql"
	"math"
	"time"

	"nescloud/backend-app/internal/apps/domain/entity"
	"nescloud/backend-app/internal/apps/domain/repository"
	"nescloud/backend-app/internal/apps/storage"
	"nescloud/backend-app/internal/dto"
	"nescloud/backend-app/internal/errs"

	"github.com/google/uuid"
)

type Service struct {
	DB         *sql.DB
	FolderRepo *repository.FolderRepo
	FileRepo   *repository.FileRepo
	Storage    *storage.Storage
}

func NewService(db *sql.DB, folderRepo *repository.FolderRepo, fileRepo *repository.FileRepo, store *storage.Storage) *Service {
	return &Service{
		DB:         db,
		FolderRepo: folderRepo,
		FileRepo:   fileRepo,
		Storage:    store,
	}
}

func (s *Service) Create(input *dto.InputCreateFolder) (*dto.ResultFolder, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()
	folderID := uuid.NewString()

	var parentID *string
	if input.ParentFolderID != "" {
		parentID = &input.ParentFolderID
	}

	folder := &entity.Folder{
		ID:             folderID,
		OwnerID:        input.OwnerID,
		ParentFolderID: parentID,
		Name:           input.Name,
		DeletedAt:      nil,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.FolderRepo.Insert(tx, input.Ctx, folder); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.ResultFolder{
		ID:             folderID,
		Name:           input.Name,
		ParentFolderID: input.ParentFolderID,
		CreatedAt:      now.Format(time.RFC3339),
		UpdatedAt:      now.Format(time.RFC3339),
	}, nil
}

func (s *Service) List(input *dto.InputListFolders) (*dto.ResultListFolders, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	page := input.Page
	if page < 1 {
		page = 1
	}
	limit := input.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var parentID *string
	if input.ParentFolderID != "" {
		parentID = &input.ParentFolderID
	}

	folders, total, err := s.FolderRepo.FindByOwnerID(tx, input.Ctx, input.OwnerID, parentID, limit, offset)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	var resultFolders []dto.ResultFolder
	for _, f := range folders {
		var parentID string
		if f.ParentFolderID != nil {
			parentID = *f.ParentFolderID
		}
		resultFolders = append(resultFolders, dto.ResultFolder{
			ID:             f.ID,
			Name:           f.Name,
			ParentFolderID: parentID,
			CreatedAt:      f.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      f.UpdatedAt.Format(time.RFC3339),
		})
	}

	if resultFolders == nil {
		resultFolders = []dto.ResultFolder{}
	}

	return &dto.ResultListFolders{
		Folders: resultFolders,
		Pagination: dto.ResultPagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) Update(input *dto.InputUpdateFolder) (*dto.ResultFolder, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	folder, err := s.FolderRepo.FindByID(tx, input.Ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if folder.OwnerID != input.OwnerID {
		return nil, errs.ErrDataNotFound
	}

	if err := s.FolderRepo.UpdateName(tx, input.Ctx, input.ID, input.Name); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.ResultFolder{
		ID:   input.ID,
		Name: input.Name,
	}, nil
}

func (s *Service) ListTrash(input *dto.InputListTrash) ([]dto.ResultTrashItem, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	folders, err := s.FolderRepo.FindDeletedByOwnerID(tx, input.Ctx, input.OwnerID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var result []dto.ResultTrashItem
	for _, f := range folders {
		var parentID string
		if f.ParentFolderID != nil {
			parentID = *f.ParentFolderID
		}
		var deletedAt string
		if f.DeletedAt != nil {
			deletedAt = f.DeletedAt.Format(time.RFC3339)
		}
		result = append(result, dto.ResultTrashItem{
			ID:             f.ID,
			OwnerID:        f.OwnerID,
			ParentFolderID: parentID,
			Name:           f.Name,
			DeletedAt:      deletedAt,
			CreatedAt:      f.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      f.UpdatedAt.Format(time.RFC3339),
		})
	}

	if result == nil {
		result = []dto.ResultTrashItem{}
	}

	return result, nil
}

func (s *Service) Delete(input *dto.InputDeleteFolder) error {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	folder, err := s.FolderRepo.FindByID(tx, input.Ctx, input.ID)
	if err != nil {
		return err
	}

	if folder.OwnerID != input.OwnerID {
		return errs.ErrDataNotFound
	}

	subfolderIDs, err := s.FolderRepo.FindSubfolderIDsRecursive(tx, input.Ctx, input.ID)
	if err != nil {
		return err
	}

	if len(subfolderIDs) > 0 {
		filesInTree, err := s.FileRepo.FindFilesByFolderIDs(tx, input.Ctx, subfolderIDs)
		if err != nil {
			return err
		}

		if len(filesInTree) > 0 {
			var fileIDs []string
			for _, f := range filesInTree {
				fileIDs = append(fileIDs, f.ID)
			}
			if err := s.FileRepo.BulkSoftDelete(tx, input.Ctx, fileIDs, time.Now()); err != nil {
				return err
			}
		}

		for _, id := range subfolderIDs {
			if err := s.FolderRepo.SoftDelete(tx, input.Ctx, id, time.Now()); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) Restore(input *dto.InputRestoreFolder) error {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	folder, err := s.FolderRepo.FindByID(tx, input.Ctx, input.ID)
	if err != nil {
		return err
	}

	if folder.OwnerID != input.OwnerID {
		return errs.ErrDataNotFound
	}

	if folder.DeletedAt == nil {
		return errs.ErrDataNotFound
	}

	if err := s.FolderRepo.Restore(tx, input.Ctx, input.ID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) PermanentDelete(input *dto.InputPermanentDeleteFolder) error {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	folder, err := s.FolderRepo.FindByID(tx, input.Ctx, input.ID)
	if err != nil {
		return err
	}

	if folder.OwnerID != input.OwnerID {
		return errs.ErrDataNotFound
	}

	subfolderIDs, err := s.FolderRepo.FindSubfolderIDsRecursive(tx, input.Ctx, input.ID)
	if err != nil {
		return err
	}

	if len(subfolderIDs) > 0 {
		filesInTree, err := s.FileRepo.FindFilesByFolderIDs(tx, input.Ctx, subfolderIDs)
		if err != nil {
			return err
		}

		for _, f := range filesInTree {
			_ = s.Storage.DeleteObject(input.Ctx, f.StorageKey)
		}

		if len(filesInTree) > 0 {
			var fileIDs []string
			for _, f := range filesInTree {
				fileIDs = append(fileIDs, f.ID)
			}
			if err := s.FileRepo.BulkHardDelete(tx, input.Ctx, fileIDs); err != nil {
				return err
			}
		}

		if err := s.FolderRepo.BulkHardDelete(tx, input.Ctx, subfolderIDs); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) EmptyTrash(input *dto.InputEmptyTrash) error {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	deletedFolders, err := s.FolderRepo.FindDeletedByOwnerID(tx, input.Ctx, input.OwnerID)
	if err != nil {
		return err
	}

	deletedFiles, err := s.FileRepo.FindDeletedByOwnerID(tx, input.Ctx, input.OwnerID)
	if err != nil {
		return err
	}

	for _, f := range deletedFiles {
		_ = s.Storage.DeleteObject(input.Ctx, f.StorageKey)
	}

	if len(deletedFiles) > 0 {
		var fileIDs []string
		for _, f := range deletedFiles {
			fileIDs = append(fileIDs, f.ID)
		}
		if err := s.FileRepo.BulkHardDelete(tx, input.Ctx, fileIDs); err != nil {
			return err
		}
	}

	if len(deletedFolders) > 0 {
		var folderIDs []string
		for _, f := range deletedFolders {
			folderIDs = append(folderIDs, f.ID)
		}
		if err := s.FolderRepo.BulkHardDelete(tx, input.Ctx, folderIDs); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
