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

type InputListTrash struct {
	Ctx     context.Context
	OwnerID string
}

type ResultTrashItem struct {
	ID             string `json:"id"`
	OwnerID        string `json:"owner_id"`
	ParentFolderID string `json:"parent_folder_id"`
	Name           string `json:"name"`
	DeletedAt      string `json:"deleted_at"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
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

type InputInitUpload struct {
	Ctx      context.Context
	OwnerID  string
	FolderID string
	FileName string
	MimeType string
	Size     int64
}

type InputCompleteUpload struct {
	Ctx    context.Context
	ID     string
	UserID string
}

type InputListFiles struct {
	Ctx      context.Context
	OwnerID  string
	FolderID string
	Page     int
	Limit    int
	Search   string
}

type InputRenameFile struct {
	Ctx     context.Context
	ID      string
	OwnerID string
	Name    string
}

type InputSoftDeleteFile struct {
	Ctx     context.Context
	ID      string
	OwnerID string
}

type InputRestoreFile struct {
	Ctx     context.Context
	ID      string
	OwnerID string
}

type InputPermanentDeleteFile struct {
	Ctx     context.Context
	ID      string
	OwnerID string
}

type InputListTrashFiles struct {
	Ctx     context.Context
	OwnerID string
}

type InputDownloadFile struct {
	Ctx     context.Context
	ID      string
	OwnerID string
}

type InputEmptyTrash struct {
	Ctx     context.Context
	OwnerID string
}

type InputGeneratePublicLink struct {
	Ctx     context.Context
	FileID  string
	OwnerID string
	BaseURL string
}

type InputRevokePublicLink struct {
	Ctx     context.Context
	FileID  string
	OwnerID string
}

type InputAccessPublicLink struct {
	Ctx   context.Context
	Token string
}

type InputListPublicLinks struct {
	Ctx     context.Context
	OwnerID string
	BaseURL string
}

type ResultFile struct {
	ID           string `json:"id"`
	OwnerID      string `json:"owner_id"`
	FolderID     string `json:"folder_id"`
	Name         string `json:"name"`
	MimeType     string `json:"mime_type"`
	Extension    string `json:"extension"`
	Size         int64  `json:"size"`
	UploadStatus string `json:"upload_status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ResultListFiles struct {
	Files      []ResultFile
	Pagination ResultPagination
}

type ResultTrashFileItem struct {
	ID           string `json:"id"`
	OwnerID      string `json:"owner_id"`
	FolderID     string `json:"folder_id"`
	Name         string `json:"name"`
	MimeType     string `json:"mime_type"`
	Extension    string `json:"extension"`
	Size         int64  `json:"size"`
	UploadStatus string `json:"upload_status"`
	DeletedAt    string `json:"deleted_at"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
