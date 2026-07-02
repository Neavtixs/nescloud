package publiclink

import (
	"database/sql"
	"errors"
	"time"

	"nescloud/backend-app/internal/apps/domain/entity"
	"nescloud/backend-app/internal/apps/domain/repository"
	"nescloud/backend-app/internal/apps/storage"
	"nescloud/backend-app/internal/dto"
	"nescloud/backend-app/internal/errs"

	"github.com/google/uuid"
)

type Service struct {
	DB             *sql.DB
	PublicLinkRepo *repository.PublicLinkRepo
	FileRepo       *repository.FileRepo
	Storage        *storage.Storage
}

func NewService(db *sql.DB, publicLinkRepo *repository.PublicLinkRepo, fileRepo *repository.FileRepo, store *storage.Storage) *Service {
	return &Service{
		DB:             db,
		PublicLinkRepo: publicLinkRepo,
		FileRepo:       fileRepo,
		Storage:        store,
	}
}

func (s *Service) Generate(input *dto.InputGeneratePublicLink) (*dto.PublicLinkRes, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	file, err := s.FileRepo.FindByID(tx, input.Ctx, input.FileID)
	if err != nil {
		return nil, err
	}

	if file.OwnerID != input.OwnerID {
		return nil, errs.ErrDataNotFound
	}

	existing, err := s.PublicLinkRepo.FindByFileID(tx, input.Ctx, input.FileID)
	if err == nil && existing != nil {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &dto.PublicLinkRes{
			URL:   input.BaseURL + "/public/" + existing.Token,
			Token: existing.Token,
		}, nil
	}

	link := &entity.PublicLink{
		ID:        uuid.NewString(),
		FileID:    input.FileID,
		Token:     uuid.NewString(),
		ExpiredAt: nil,
		CreatedAt: time.Now(),
	}

	if err := s.PublicLinkRepo.Insert(tx, input.Ctx, link); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.PublicLinkRes{
		URL:   input.BaseURL + "/public/" + link.Token,
		Token: link.Token,
	}, nil
}

func (s *Service) Revoke(input *dto.InputRevokePublicLink) error {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	file, err := s.FileRepo.FindByID(tx, input.Ctx, input.FileID)
	if err != nil {
		return err
	}

	if file.OwnerID != input.OwnerID {
		return errs.ErrDataNotFound
	}

	if err := s.PublicLinkRepo.Delete(tx, input.Ctx, input.FileID); err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			if err := tx.Commit(); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) Access(input *dto.InputAccessPublicLink) (string, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	link, err := s.PublicLinkRepo.FindByToken(tx, input.Ctx, input.Token)
	if err != nil {
		return "", err
	}

	if link.ExpiredAt != nil && time.Now().After(*link.ExpiredAt) {
		return "", errs.ErrDataNotFound
	}

	file, err := s.FileRepo.FindByID(tx, input.Ctx, link.FileID)
	if err != nil {
		return "", err
	}

	if file.DeletedAt != nil {
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

func (s *Service) ListLinks(input *dto.InputListPublicLinks) ([]dto.PublicLinkItemRes, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	items, err := s.PublicLinkRepo.FindAllByOwnerID(tx, input.Ctx, input.OwnerID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var result []dto.PublicLinkItemRes
	for _, item := range items {
		result = append(result, dto.PublicLinkItemRes{
			FileID:    item.FileID,
			FileName:  item.OriginalName,
			MimeType:  item.MimeType,
			Size:      item.Size,
			PublicURL: input.BaseURL + "/public/" + item.Token,
			CreatedAt: item.CreatedAt.Format(time.RFC3339),
		})
	}
	if result == nil {
		result = []dto.PublicLinkItemRes{}
	}

	return result, nil
}
