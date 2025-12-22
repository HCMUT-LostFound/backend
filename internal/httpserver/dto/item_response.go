package dto

import (
	"time"

	"github.com/google/uuid"
)

type ItemResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"userId"`

	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`

	ImageURLs   []string   `json:"imageUrls"`

	Location    string     `json:"location"`
	Campus      string     `json:"campus"`

	LostAt      *time.Time `json:"lostAt"`

	Tags        []string   `json:"tags"`

	CreatedAt   time.Time  `json:"createdAt"`
}
