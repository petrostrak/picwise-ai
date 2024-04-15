package types

import (
	"github.com/google/uuid"
	"time"
)

type ImageStatus int

const (
	ImageStatusFailed ImageStatus = iota
	ImageStatusPending
	ImageStatusCompleted
)

type Image struct {
	ID            int `bun:"id,pk,autoincrement"`
	UserID        uuid.UUID
	Status        ImageStatus
	Prompt        string
	ImageLocation string
	Deleted       bool      `bun:"default:'false'"`
	CreatedAt     time.Time `bun:"default:'now()'"`
	DeletedAt     time.Time
}
