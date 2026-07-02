package entity

import "time"

type PublicLink struct {
	ID        string
	FileID    string
	Token     string
	ExpiredAt *time.Time
	CreatedAt time.Time
}

type PublicLinkFile struct {
	ID           string
	FileID       string
	Token        string
	ExpiredAt    *time.Time
	CreatedAt    time.Time
	OriginalName string
	MimeType     string
	Size         int64
}
