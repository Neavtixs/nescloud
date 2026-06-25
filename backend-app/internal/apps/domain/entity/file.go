package entity

import "time"

type File struct {
	ID           string
	OwnerID      string
	FolderID     string
	OriginalName string
	StorageKey   string
	MimeType     string
	Extension    string
	Size         int64
	Checksum     string
	UploadStatus string
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
