package dto

type ResponseWeb[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type AuthRegisterReq struct {
	Name     string `json:"name"     validate:"required,min=1,max=100"`
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthLoginReq struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthRegisterRes struct {
	AccessToken string `json:"access_token"`
}

type AuthLoginRes struct {
	AccessToken string `json:"access_token"`
}

type AuthMeRes struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type AuthRefreshRes struct {
	AccessToken string `json:"access_token"`
}

type FolderCreateReq struct {
	ParentFolderID string `json:"parent_folder_id"`
	Name           string `json:"name" validate:"required,min=1,max=255"`
}

type FolderUpdateReq struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type FolderRes struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ParentFolderID string `json:"parent_folder_id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type FolderIDRes struct {
	ID string `json:"id"`
}

type PaginationRes struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type FileCompleteReq struct {
	FileID string `json:"file_id" validate:"required"`
}

type FileInitUploadReq struct {
	FolderID string `json:"folder_id"`
	FileName string `json:"file_name" validate:"required,min=1,max=255"`
	MimeType string `json:"mime_type" validate:"required"`
	Size     int64  `json:"size" validate:"required"`
}

type FileRenameReq struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type FileRes struct {
	ID           string `json:"id"`
	FolderID     string `json:"folder_id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	Extension    string `json:"extension"`
	Size         int64  `json:"size"`
	UploadStatus string `json:"upload_status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type FileUploadRes struct {
	FileID    string `json:"file_id"`
	UploadURL string `json:"upload_url"`
	ExpiresAt string `json:"expired_at"`
}

type PublicLinkRes struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type PublicLinkItemRes struct {
	FileID    string `json:"file_id"`
	FileName  string `json:"file_name"`
	MimeType  string `json:"mime_type"`
	Size      int64  `json:"size"`
	PublicURL string `json:"public_url"`
	CreatedAt string `json:"created_at"`
}

type ErrorWeb struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}
