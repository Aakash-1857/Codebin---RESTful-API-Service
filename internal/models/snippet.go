package models

import (
	"errors"
	"time"
)
var ErrNoRecord=errors.New("models:no matching record found")
// In internal/models/snippet.go

type Snippet struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    ExpiresAt time.Time `json:"expires_at"`
}