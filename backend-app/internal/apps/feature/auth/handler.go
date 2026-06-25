package auth

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

func (h *Handler) RegisterHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)

	var req dto.AuthRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
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

	input := &dto.InputAuthRegister{
		Ctx:      c.Request.Context(),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.Service.Register(input)
	if err != nil {
		if errors.Is(err, errs.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "auth_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	c.SetCookie("refresh_token", result.RefreshToken, 604800, "/api/auth", "", false, true)

	c.JSON(http.StatusCreated, dto.ResponseWeb[*dto.ResultAuthRegister]{
		Message: "register user success",
		Data:    result,
	})
}
