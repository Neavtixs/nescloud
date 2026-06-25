package entity

import "time"

type Quota struct {
	ID          string
	UserID      string
	MaxStorage  int64
	UsedStorage int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
