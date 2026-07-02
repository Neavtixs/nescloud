package file

import (
	"database/sql"
	"math"
	"path/filepath"
	"time"

	"nescloud/backend-app/internal/apps/domain/entity"
	"nescloud/backend-app/internal/apps/domain/repository"
	"nescloud/backend-app/internal/apps/storage"
	"nescloud/backend-app/internal/dto"
	"nescloud/backend-app/internal/errs"

	"github.com/google/uuid"
)

type Service struct {
	DB        *sql.DB
	FileRepo  *repository.FileRepo
	QuotaRepo *repository.QuotaRepo
	Storage   *storage.Storage
}

func NewService(db *sql.DB, fileRepo *repository.FileRepo, quotaRepo *repository.QuotaRepo, store *storage.Storage) *Service {
	return &Service{
		DB:        db,
		FileRepo:  fileRepo,
		QuotaRepo: quotaRepo,
		Storage:   store,
	}
}

func (s *Service) InitUpload(input *dto.InputInitUpload) (*dto.ResultFile, string, time.Time, error) {
	const maxSize int64 = 100 * 1024 * 1024

	if input.Size > maxSize {
		return nil, "", time.Time{}, errs.ErrFileTooLarge
	}

	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, "", time.Time{}, err
	}
	defer tx.Rollback()

	now := time.Now()
	fileID := uuid.NewString()
	ext := filepath.Ext(input.FileName)
	storageKey := input.OwnerID + "/" + fileID + ext

	var folderID *string
	if input.FolderID != "" {
		folderID = &input.FolderID
	}

	file := &entity.File{
		ID:           fileID,
		OwnerID:      input.OwnerID,
		FolderID:     folderID,
		OriginalName: input.FileName,
		StorageKey:   storageKey,
		MimeType:     input.MimeType,
		Extension:    ext,
		Size:         input.Size,
		Checksum:     "",
		UploadStatus: "uploading",
		DeletedAt:    nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.FileRepo.Insert(tx, input.Ctx, file); err != nil {
		return nil, "", time.Time{}, err
	}

	if err := tx.Commit(); err != nil {
		return nil, "", time.Time{}, err
	}

	uploadURL, expiresAt, err := s.Storage.GenerateUploadURL(input.Ctx, storageKey, 15*time.Minute)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	result := &dto.ResultFile{
		ID:           fileID,
		OwnerID:      input.OwnerID,
		FolderID:     input.FolderID,
		Name:         input.FileName,
		MimeType:     input.MimeType,
		Extension:    ext,
		Size:         input.Size,
		UploadStatus: "uploading",
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
	}

	return result, uploadURL, expiresAt, nil
}

func (s *Service) CompleteUpload(input *dto.InputCompleteUpload) error {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	file, err := s.FileRepo.FindByID(tx, input.Ctx, input.ID)
	if err != nil {
		return err
	}

	if file.OwnerID != input.UserID {
		return errs.ErrDataNotFound
	}

	exists, err := s.Storage.ObjectExists(input.Ctx, file.StorageKey)
	if err != nil || !exists {
		return errs.ErrDataNotFound
	}

	if err := s.FileRepo.UpdateStatus(tx, input.Ctx, input.ID, "completed"); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) ListFiles(input *dto.InputListFiles) (*dto.ResultListFiles, error) {
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

	var folderID *string
	if input.FolderID != "" {
		folderID = &input.FolderID
	}

	files, total, err := s.FileRepo.FindByOwnerIDAndFolderID(tx, input.Ctx, input.OwnerID, folderID, input.Search, limit, offset)
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

	var resultFiles []dto.ResultFile
	for _, f := range files {
		var folderIDStr string
		if f.FolderID != nil {
			folderIDStr = *f.FolderID
		}
		resultFiles = append(resultFiles, dto.ResultFile{
			ID:           f.ID,
			OwnerID:      f.OwnerID,
			FolderID:     folderIDStr,
			Name:         f.OriginalName,
			MimeType:     f.MimeType,
			Extension:    f.Extension,
			Size:         f.Size,
			UploadStatus: f.UploadStatus,
			CreatedAt:    f.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    f.UpdatedAt.Format(time.RFC3339),
		})
	}

	if resultFiles == nil {
		resultFiles = []dto.ResultFile{}
	}

	return &dto.ResultListFiles{
		Files: resultFiles,
		Pagination: dto.ResultPagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) Rename(input *dto.InputRenameFile) (*dto.ResultFile, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	file, err := s.FileRepo.FindByID(tx, input.Ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if file.OwnerID != input.OwnerID {
		return nil, errs.ErrDataNotFound
	}

	if err := s.FileRepo.UpdateName(tx, input.Ctx, input.ID, input.Name); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.ResultFile{
		ID:      input.ID,
		Name:    input.Name,
		OwnerID: input.OwnerID,
	}, nil
}

func (s *Service) Delete(input *dto.InputSoftDeleteFile) error {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	file, err := s.FileRepo.FindByID(tx, input.Ctx, input.ID)
	if err != nil {
		return err
	}

	if file.OwnerID != input.OwnerID {
		return errs.ErrDataNotFound
	}

	if err := s.FileRepo.SoftDelete(tx, input.Ctx, input.ID, time.Now()); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) Restore(input *dto.InputRestoreFile) error {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	file, err := s.FileRepo.FindByID(tx, input.Ctx, input.ID)
	if err != nil {
		return err
	}

	if file.OwnerID != input.OwnerID {
		return errs.ErrDataNotFound
	}

	if file.DeletedAt == nil {
		return errs.ErrDataNotFound
	}

	if err := s.FileRepo.Restore(tx, input.Ctx, input.ID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) PermanentDelete(input *dto.InputPermanentDeleteFile) error {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	file, err := s.FileRepo.FindByID(tx, input.Ctx, input.ID)
	if err != nil {
		return err
	}

	if file.OwnerID != input.OwnerID {
		return errs.ErrDataNotFound
	}

	if err := s.Storage.DeleteObject(input.Ctx, file.StorageKey); err != nil {
		return err
	}

	if err := s.FileRepo.HardDelete(tx, input.Ctx, input.ID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) ListTrash(input *dto.InputListTrashFiles) ([]dto.ResultTrashFileItem, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	files, err := s.FileRepo.FindDeletedByOwnerID(tx, input.Ctx, input.OwnerID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var result []dto.ResultTrashFileItem
	for _, f := range files {
		var deletedAt string
		if f.DeletedAt != nil {
			deletedAt = f.DeletedAt.Format(time.RFC3339)
		}
		var folderIDStr string
		if f.FolderID != nil {
			folderIDStr = *f.FolderID
		}
		result = append(result, dto.ResultTrashFileItem{
			ID:           f.ID,
			OwnerID:      f.OwnerID,
			FolderID:     folderIDStr,
			Name:         f.OriginalName,
			MimeType:     f.MimeType,
			Extension:    f.Extension,
			Size:         f.Size,
			UploadStatus: f.UploadStatus,
			DeletedAt:    deletedAt,
			CreatedAt:    f.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    f.UpdatedAt.Format(time.RFC3339),
		})
	}

	if result == nil {
		result = []dto.ResultTrashFileItem{}
	}

	return result, nil
}

func (s *Service) GetFile(input *dto.InputDownloadFile) (*dto.ResultFile, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	file, err := s.FileRepo.FindByID(tx, input.Ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if file.OwnerID != input.OwnerID {
		return nil, errs.ErrDataNotFound
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var folderIDStr string
	if file.FolderID != nil {
		folderIDStr = *file.FolderID
	}

	return &dto.ResultFile{
		ID:           file.ID,
		OwnerID:      file.OwnerID,
		FolderID:     folderIDStr,
		Name:         file.OriginalName,
		MimeType:     file.MimeType,
		Extension:    file.Extension,
		Size:         file.Size,
		UploadStatus: file.UploadStatus,
		CreatedAt:    file.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    file.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Service) GetDownloadURL(input *dto.InputDownloadFile) (string, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	file, err := s.FileRepo.FindByID(tx, input.Ctx, input.ID)
	if err != nil {
		return "", err
	}

	if file.OwnerID != input.OwnerID {
		return "", errs.ErrDataNotFound
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	downloadURL, err := s.Storage.GenerateDownloadURL(input.Ctx, file.StorageKey, 15*time.Minute)
	if err != nil {
		return "", err
	}

	return downloadURL, nil
}
