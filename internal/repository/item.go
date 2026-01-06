package repository

import (
	"time"
	"github.com/lib/pq"
	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID `db:"id"`
	UserID string `db:"user_id"`

	Type        string    `db:"type"` // lost | found
	Title       string    `db:"title"`
	Description *string   `db:"description"`

	ImageURLs   pq.StringArray  `db:"image_urls"`

	Location    string    `db:"location"`
	Campus      string    `db:"campus"`

	LostAt      *time.Time `db:"lost_at"`

	Tags        pq.StringArray   `db:"tags"`

	IsConfirmed bool       `db:"is_confirmed"`

	CreatedAt   time.Time  `db:"created_at"`
}
