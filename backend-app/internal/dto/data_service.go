package dto

import "context"

type InputAuthRegister struct {
	Ctx      context.Context
	Name     string
	Email    string
	Password string
}

type InputAuthLogin struct {
	Ctx      context.Context
	Email    string
	Password string
}

type InputAuthRefresh struct {
	Ctx          context.Context
	RefreshToken string
}

type ResultAuthLogin struct {
	AccessToken      string
	RefreshToken     string
	RefreshExpiresIn int
}

type ResultAuthMe struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type ResultAuthRegister struct {
	ID               string
	AccessToken      string
	RefreshToken     string
	RefreshExpiresIn int
}

type InputAuthLogout struct {
	Ctx          context.Context
	RefreshToken string
}

type InputAuthMe struct {
	Ctx    context.Context
	UserID string
}

type ResultAuthRefresh struct {
	AccessToken      string
	RefreshToken     string
	RefreshExpiresIn int
}

type InputCreateFolder struct {
	Ctx            context.Context
	OwnerID        string
	ParentFolderID string
	Name           string
}

type InputListFolders struct {
	Ctx            context.Context
	OwnerID        string
	ParentFolderID string
	Page           int
	Limit          int
}

type InputUpdateFolder struct {
	Ctx     context.Context
	ID      string
	OwnerID string
	Name    string
}

type InputDeleteFolder struct {
	Ctx     context.Context
	ID      string
	OwnerID string
}

type InputRestoreFolder struct {
	Ctx     context.Context
	ID      string
	OwnerID string
}

type InputPermanentDeleteFolder struct {
	Ctx     context.Context
	ID      string
	OwnerID string
}

type ResultFolder struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ParentFolderID string `json:"parent_folder_id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type ResultListFolders struct {
	Folders    []ResultFolder
	Pagination ResultPagination
}

type ResultPagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
