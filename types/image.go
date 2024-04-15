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
	ImageLocation string
	CreatedAt     time.Time `bun:"default:'now()'"`
}
