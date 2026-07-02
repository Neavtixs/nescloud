package folder

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

func (h *Handler) CreateFolderHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("create folder request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	var req dto.FolderCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithField("layer", "folder_handler").Warn("invalid request format")
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

	input := &dto.InputCreateFolder{
		Ctx:            c.Request.Context(),
		OwnerID:        userIDStr,
		ParentFolderID: req.ParentFolderID,
		Name:           req.Name,
	}

	result, err := h.Service.Create(input)
	if err != nil {
		log.WithField("layer", "folder_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("folder_id", result.ID).Info("folder created")

	c.JSON(http.StatusCreated, dto.ResponseWeb[dto.FolderIDRes]{
		Message: "folder created",
		Data: dto.FolderIDRes{
			ID: result.ID,
		},
	})
}

func (h *Handler) ListFoldersHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("list folders request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	parentFolderID := c.Query("parent_folder_id")
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

	input := &dto.InputListFolders{
		Ctx:            c.Request.Context(),
		OwnerID:        userIDStr,
		ParentFolderID: parentFolderID,
		Page:           page,
		Limit:          limit,
	}

	result, err := h.Service.List(input)
	if err != nil {
		log.WithField("layer", "folder_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.Info("folders listed")

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    result.Folders,
		"pagination": dto.PaginationRes{
			Page:       result.Pagination.Page,
			Limit:      result.Pagination.Limit,
			Total:      result.Pagination.Total,
			TotalPages: result.Pagination.TotalPages,
		},
	})
}

func (h *Handler) UpdateFolderHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("update folder request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	folderID := c.Param("id")
	if folderID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	var req dto.FolderUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithField("layer", "folder_handler").Warn("invalid request format")
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

	input := &dto.InputUpdateFolder{
		Ctx:     c.Request.Context(),
		ID:      folderID,
		OwnerID: userIDStr,
		Name:    req.Name,
	}

	result, err := h.Service.Update(input)
	if err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.WithField("folder_id", folderID).Warn("folder not found")
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "folder_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("folder_id", result.ID).Info("folder renamed")

	c.JSON(http.StatusOK, dto.ResponseWeb[dto.FolderRes]{
		Message: "folder renamed",
		Data: dto.FolderRes{
			ID:   result.ID,
			Name: result.Name,
		},
	})
}

func (h *Handler) ListTrashHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("list trash request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	input := &dto.InputListTrash{
		Ctx:     c.Request.Context(),
		OwnerID: userIDStr,
	}

	result, err := h.Service.ListTrash(input)
	if err != nil {
		log.WithField("layer", "folder_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.Info("trash items listed")

	c.JSON(http.StatusOK, dto.ResponseWeb[[]dto.ResultTrashItem]{
		Message: "success",
		Data:    result,
	})
}

func (h *Handler) DeleteFolderHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("delete folder request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	folderID := c.Param("id")
	if folderID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	input := &dto.InputDeleteFolder{
		Ctx:     c.Request.Context(),
		ID:      folderID,
		OwnerID: userIDStr,
	}

	if err := h.Service.Delete(input); err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.WithField("folder_id", folderID).Warn("folder not found")
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "folder_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("folder_id", folderID).Info("folder moved to trash")

	c.JSON(http.StatusOK, dto.ResponseWeb[any]{
		Message: "folder moved to trash",
	})
}

func (h *Handler) RestoreFolderHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("restore folder request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	folderID := c.Param("id")
	if folderID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	input := &dto.InputRestoreFolder{
		Ctx:     c.Request.Context(),
		ID:      folderID,
		OwnerID: userIDStr,
	}

	if err := h.Service.Restore(input); err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.WithField("folder_id", folderID).Warn("folder not found")
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "folder_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("folder_id", folderID).Info("folder restored")

	c.JSON(http.StatusOK, dto.ResponseWeb[any]{
		Message: "folder restored",
	})
}

func (h *Handler) PermanentDeleteFolderHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("permanent delete folder request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	folderID := c.Param("id")
	if folderID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	input := &dto.InputPermanentDeleteFolder{
		Ctx:     c.Request.Context(),
		ID:      folderID,
		OwnerID: userIDStr,
	}

	if err := h.Service.PermanentDelete(input); err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.WithField("folder_id", folderID).Warn("folder not found")
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "folder_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("folder_id", folderID).Info("folder permanently deleted")

	c.JSON(http.StatusOK, dto.ResponseWeb[any]{
		Message: "folder permanently deleted",
	})
}

func (h *Handler) EmptyTrashHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("empty trash request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	input := &dto.InputEmptyTrash{
		Ctx:     c.Request.Context(),
		OwnerID: userIDStr,
	}

	if err := h.Service.EmptyTrash(input); err != nil {
		log.WithField("layer", "folder_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.Info("trash emptied")

	c.JSON(http.StatusOK, dto.ResponseWeb[any]{
		Message: "trash emptied",
	})
}
