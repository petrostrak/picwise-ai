package types

import (
	"github.com/google/uuid"
	"time"
)

type Account struct {
	ID       int `bun:"id,pk,autoincrement"`
	UserID   uuid.UUID
	Username string
	CreateAt time.Time `bun:"default:'now()'"`
}
