package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	UUID      uuid.UUID `json:"uuid"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
