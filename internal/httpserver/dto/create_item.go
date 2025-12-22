package dto

import "time"

type CreateItemRequest struct {
	Type        string     `json:"type" binding:"required"`
	Title       string     `json:"title" binding:"required"`
	Description *string    `json:"description"`

	ImageURLs   []string   `json:"imageUrls"`

	Location    string     `json:"location" binding:"required"`
	Campus      string     `json:"campus" binding:"required"`

	LostAt      *time.Time `json:"lostAt"`

	Tags        []string   `json:"tags"`
}
