package publiclink

import (
	"errors"
	"net/http"

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

func (h *Handler) GenerateHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("generate public link request received")

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

	scheme := "https"
	if c.Request.TLS == nil {
		scheme = "http"
	}
	baseURL := scheme + "://" + c.Request.Host

	input := &dto.InputGeneratePublicLink{
		Ctx:     c.Request.Context(),
		FileID:  fileID,
		OwnerID: userIDStr,
		BaseURL: baseURL,
	}

	result, err := h.Service.Generate(input)
	if err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.WithField("file_id", fileID).Warn("file not found")
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "publiclink_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("file_id", fileID).Info("public link generated")

	c.JSON(http.StatusCreated, dto.ResponseWeb[dto.PublicLinkRes]{
		Message: "public link created",
		Data:    *result,
	})
}

func (h *Handler) RevokeHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("revoke public link request received")

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

	input := &dto.InputRevokePublicLink{
		Ctx:     c.Request.Context(),
		FileID:  fileID,
		OwnerID: userIDStr,
	}

	if err := h.Service.Revoke(input); err != nil {
		log.WithField("layer", "publiclink_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("file_id", fileID).Info("public link revoked")

	c.JSON(http.StatusOK, dto.ResponseWeb[any]{
		Message: "public link revoked",
	})
}

func (h *Handler) AccessHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("access public file request received")

	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorWeb{Message: "invalid request format"})
		return
	}

	input := &dto.InputAccessPublicLink{
		Ctx:   c.Request.Context(),
		Token: token,
	}

	downloadURL, err := h.Service.Access(input)
	if err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorWeb{Message: "public link not found"})
			return
		}
		log.WithField("layer", "publiclink_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.Info("public file accessed")
	c.Redirect(http.StatusFound, downloadURL)
}

func (h *Handler) ListLinksHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("list public links request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	scheme := "https"
	if c.Request.TLS == nil {
		scheme = "http"
	}
	baseURL := scheme + "://" + c.Request.Host

	input := &dto.InputListPublicLinks{
		Ctx:     c.Request.Context(),
		OwnerID: userIDStr,
		BaseURL: baseURL,
	}

	result, err := h.Service.ListLinks(input)
	if err != nil {
		log.WithField("layer", "publiclink_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.Info("public links listed")

	c.JSON(http.StatusOK, dto.ResponseWeb[[]dto.PublicLinkItemRes]{
		Message: "success",
		Data:    result,
	})
}
