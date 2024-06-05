package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
