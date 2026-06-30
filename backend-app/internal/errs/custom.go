package errs

import "errors"

var (
	ErrInternal                = errors.New("internal server error")
	ErrDataNotFound            = errors.New("data not found")
	ErrInvalidAccessToken      = errors.New("invalid access token")
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrInvalidCredentials      = errors.New("invalid email or password")
	ErrInvalidRequest          = errors.New("invalid request format")
	ErrFolderNameAlreadyExists = errors.New("folder name already exists")
	ErrFolderHasChildren       = errors.New("folder is not empty")
)
