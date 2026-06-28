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
	log.Info("register request received")

	var req dto.AuthRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithField("layer", "auth_handler").Warn("invalid request format")
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
			log.WithField("email", req.Email).Warn("email already exists")
			c.JSON(http.StatusConflict, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "auth_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("user_id", result.ID).Info("register success")

	c.SetCookie("refresh_token", result.RefreshToken, result.RefreshExpiresIn, "/", "", false, true)

	c.JSON(http.StatusCreated, dto.ResponseWeb[dto.AuthRegisterRes]{
		Message: "register user success",
		Data: dto.AuthRegisterRes{
			AccessToken: result.AccessToken,
		},
	})
}

func (h *Handler) LoginHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("login request received")

	var req dto.AuthLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithField("layer", "auth_handler").Warn("invalid request format")
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

	input := &dto.InputAuthLogin{
		Ctx:      c.Request.Context(),
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.Service.Login(input)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			log.WithField("email", req.Email).Warn("invalid login credentials")
			c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: err.Error()})
			return
		}
		log.WithField("layer", "auth_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.Info("login success")

	c.SetCookie("refresh_token", result.RefreshToken, result.RefreshExpiresIn*2, "/", "", false, true)

	c.JSON(http.StatusOK, dto.ResponseWeb[dto.AuthLoginRes]{
		Message: "login user success",
		Data: dto.AuthLoginRes{
			AccessToken: result.AccessToken,
		},
	})
}

func (h *Handler) RefreshHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("refresh request received")

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		log.Warn("refresh_token cookie not found")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}
	log.Info("cookie: " + refreshToken)

	input := &dto.InputAuthRefresh{
		Ctx:          c.Request.Context(),
		RefreshToken: refreshToken,
	}

	log.Info("refresh service call")
	result, err := h.Service.Refresh(input)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidAccessToken) {
			log.Warn("invalid refresh token")
			c.SetCookie("refresh_token", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
			return
		}
		log.WithField("layer", "auth_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.Info("refresh success")

	c.SetCookie("refresh_token", result.RefreshToken, result.RefreshExpiresIn*2, "/", "", false, true)

	c.JSON(http.StatusOK, dto.ResponseWeb[dto.AuthRefreshRes]{
		Message: "token refreshed",
		Data: dto.AuthRefreshRes{
			AccessToken: result.AccessToken,
		},
	})
}

func (h *Handler) LogoutHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("logout request received")

	refreshToken, _ := c.Cookie("refresh_token")
	if refreshToken == "" {
		log.Warn("no refresh_token cookie found")
	} else {
		log.Info("refresh_token cookie found, proceeding with logout")
	}

	input := &dto.InputAuthLogout{
		Ctx:          c.Request.Context(),
		RefreshToken: refreshToken,
	}

	_ = h.Service.Logout(input)

	c.SetCookie("refresh_token", "", 0, "/", "", false, true)

	log.Info("logout success")

	c.JSON(http.StatusOK, dto.ResponseWeb[any]{
		Message: "logout success",
	})
}

func (h *Handler) MeHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log.Info("me request received")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		log.Warn("invalid or missing user_id in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
		return
	}

	input := &dto.InputAuthMe{
		Ctx:    c.Request.Context(),
		UserID: userIDStr,
	}

	result, err := h.Service.Me(input)
	if err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			log.Warn("user not found")
			c.JSON(http.StatusUnauthorized, dto.ErrorWeb{Message: errs.ErrInvalidAccessToken.Error()})
			return
		}
		log.WithField("layer", "auth_handler").Error(err)
		c.JSON(http.StatusInternalServerError, dto.ErrorWeb{Message: errs.ErrInternal.Error()})
		return
	}

	log.WithField("user_id", result.ID).Info("me success")

	c.JSON(http.StatusOK, dto.ResponseWeb[dto.AuthMeRes]{
		Message: "success",
		Data: dto.AuthMeRes{
			ID:        result.ID,
			Name:      result.Name,
			Email:     result.Email,
			CreatedAt: result.CreatedAt,
		},
	})
}
