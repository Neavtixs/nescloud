package entity

import "time"

type Folder struct {
	ID             string
	OwnerID        string
	ParentFolderID *string
	Name           string
	DeletedAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
