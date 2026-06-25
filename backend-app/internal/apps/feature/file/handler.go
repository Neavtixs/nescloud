package file

import (
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
