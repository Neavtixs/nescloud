package file

import (
	"errors"
	"net/http"
	"strconv"

	"nescloud/backend-app/internal/dto"
	"nescloud/backend-app/internal/errs"
	"nescloud/backend-app/internal/helper"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Service  *Service
	Validate *validator.Validate
	Log      *logrus.Logger
}

func NewHandler(service *Service, validate *validator.Validate, log *logrus.Logger) *Handler {
	return &Handler{
		Service:  service,
		Validate: validate,
		Log:      log,
	}
}

func (h *Handler) InitUploadHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("init upload request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	var req dto.FileInitUploadReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithField("layer", "file_handler").Warn("invalid request format")
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseWeb[map[string]string]{
			Message: "validation failed",
			Data:    helper.ValidationMsg(err),
		})
		return
	}

	input := &dto.InputInitUpload{
		Ctx:      c.Request.Context(),
		OwnerID:  userIDStr,
		FolderID: req.FolderID,
		FileName: req.FileName,
		MimeType: req.MimeType,
		Size:     req.Size,
	}

	result, uploadURL, expiresAt, err := h.Service.InitUpload(input)
	if err != nil {
		if errors.Is(err, errs.ErrFileTooLarge) {
			c.JSON(http.StatusRequestEntityTooLarge, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "file_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("file_id", result.ID).Info("upload initialized")

	c.JSON(http.StatusCreated, dto.ResponseWeb[dto.FileUploadRes]{
		Message: "upload initialized",
		Data: dto.FileUploadRes{
			FileID:    result.ID,
			UploadURL: uploadURL,
			ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}

func (h *Handler) CompleteUploadHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("complete upload request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	var req dto.FileCompleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithField("layer", "file_handler").Warn("invalid request format")
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseWeb[map[string]string]{
			Message: "validation failed",
			Data:    helper.ValidationMsg(err),
		})
		return
	}

	input := &dto.InputCompleteUpload{
		Ctx:    c.Request.Context(),
		ID:     req.FileID,
		UserID: userIDStr,
	}

	if err := h.Service.CompleteUpload(input); err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.WithField("file_id", req.FileID).Warn("file not found")
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "file_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("file_id", req.FileID).Info("upload completed")

	c.JSON(http.StatusOK, dto.ResponseWeb[any]{
		Message: "upload completed",
	})
}

func (h *Handler) ListFilesHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("list files request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	folderID := c.Query("folder_id")
	search := c.Query("search")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	input := &dto.InputListFiles{
		Ctx:      c.Request.Context(),
		OwnerID:  userIDStr,
		FolderID: folderID,
		Page:     page,
		Limit:    limit,
		Search:   search,
	}

	result, err := h.Service.ListFiles(input)
	if err != nil {
		log.WithField("layer", "file_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.Info("files listed")

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    result.Files,
		"pagination": dto.PaginationRes{
			Page:       result.Pagination.Page,
			Limit:      result.Pagination.Limit,
			Total:      result.Pagination.Total,
			TotalPages: result.Pagination.TotalPages,
		},
	})
}

func (h *Handler) RenameHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("rename file request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	var req dto.FileRenameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithField("layer", "file_handler").Warn("invalid request format")
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseWeb[map[string]string]{
			Message: "validation failed",
			Data:    helper.ValidationMsg(err),
		})
		return
	}

	input := &dto.InputRenameFile{
		Ctx:     c.Request.Context(),
		ID:      fileID,
		OwnerID: userIDStr,
		Name:    req.Name,
	}

	result, err := h.Service.Rename(input)
	if err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.WithField("file_id", fileID).Warn("file not found")
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "file_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("file_id", result.ID).Info("file renamed")

	c.JSON(http.StatusOK, dto.ResponseWeb[dto.FileRes]{
		Message: "file renamed",
		Data: dto.FileRes{
			ID:           result.ID,
			OriginalName: result.Name,
		},
	})
}

func (h *Handler) DeleteFileHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("delete file request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	input := &dto.InputSoftDeleteFile{
		Ctx:     c.Request.Context(),
		ID:      fileID,
		OwnerID: userIDStr,
	}

	if err := h.Service.Delete(input); err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.WithField("file_id", fileID).Warn("file not found")
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "file_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("file_id", fileID).Info("file moved to trash")

	c.JSON(http.StatusOK, dto.ResponseWeb[any]{
		Message: "file moved to trash",
	})
}

func (h *Handler) RestoreFileHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("restore file request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	input := &dto.InputRestoreFile{
		Ctx:     c.Request.Context(),
		ID:      fileID,
		OwnerID: userIDStr,
	}

	if err := h.Service.Restore(input); err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.WithField("file_id", fileID).Warn("file not found")
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "file_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("file_id", fileID).Info("file restored")

	c.JSON(http.StatusOK, dto.ResponseWeb[any]{
		Message: "file restored",
	})
}

func (h *Handler) PermanentDeleteFileHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("permanent delete file request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	input := &dto.InputPermanentDeleteFile{
		Ctx:     c.Request.Context(),
		ID:      fileID,
		OwnerID: userIDStr,
	}

	if err := h.Service.PermanentDelete(input); err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.WithField("file_id", fileID).Warn("file not found")
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "file_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("file_id", fileID).Info("file permanently deleted")

	c.JSON(http.StatusOK, dto.ResponseWeb[any]{
		Message: "file permanently deleted",
	})
}

func (h *Handler) FileDetailHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("file detail request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	input := &dto.InputDownloadFile{
		Ctx:     c.Request.Context(),
		ID:      fileID,
		OwnerID: userIDStr,
	}

	result, err := h.Service.GetFile(input)
	if err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.WithField("file_id", fileID).Warn("file not found")
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "file_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("file_id", fileID).Info("file detail retrieved")

	c.JSON(http.StatusOK, dto.ResponseWeb[dto.FileRes]{
		Message: "success",
		Data: dto.FileRes{
			ID:           result.ID,
			FolderID:     result.FolderID,
			OriginalName: result.Name,
			MimeType:     result.MimeType,
			Extension:    result.Extension,
			Size:         result.Size,
			UploadStatus: result.UploadStatus,
			CreatedAt:    result.CreatedAt,
			UpdatedAt:    result.UpdatedAt,
		},
	})
}

func (h *Handler) DownloadHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("download file request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	input := &dto.InputDownloadFile{
		Ctx:     c.Request.Context(),
		ID:      fileID,
		OwnerID: userIDStr,
	}

	downloadURL, err := h.Service.GetDownloadURL(input)
	if err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.WithField("file_id", fileID).Warn("file not found")
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "file_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("file_id", fileID).Info("download url generated")

	c.JSON(http.StatusOK, dto.ResponseWeb[map[string]string]{
		Message: "success",
		Data:    map[string]string{"download_url": downloadURL},
	})
}

func (h *Handler) ListTrashFilesHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("list trash files request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	input := &dto.InputListTrashFiles{
		Ctx:     c.Request.Context(),
		OwnerID: userIDStr,
	}

	result, err := h.Service.ListTrash(input)
	if err != nil {
		log.WithField("layer", "file_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.Info("trash files listed")

	c.JSON(http.StatusOK, dto.ResponseWeb[[]dto.ResultTrashFileItem]{
		Message: "success",
		Data:    result,
	})
}
