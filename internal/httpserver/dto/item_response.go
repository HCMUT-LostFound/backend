package dto

import (
	"time"

	"github.com/google/uuid"
)

type Reporter struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	ImageURL string `json:"imageUrl,omitempty"`
}

type ItemResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID string `json:"userId"`
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`

	ImageURLs   []string   `json:"imageUrls"`

	Location    string     `json:"location"`
	Campus      string     `json:"campus"`

	LostAt      *time.Time `json:"lostAt"`

	Tags        []string   `json:"tags"`
	Reporter *Reporter `json:"reporter,omitempty"`
	IsConfirmed bool `json:"isConfirmed"`
	CreatedAt   time.Time  `json:"createdAt"`
}
