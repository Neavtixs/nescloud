package entity

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID           string
	UserID       string
	Action       string
	ResourceType string
	ResourceID   string
	Metadata     json.RawMessage
	CreatedAt    time.Time
}
