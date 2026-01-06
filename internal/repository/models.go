package repository

import "time"
import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID    `db:"id"`
	ClerkUserID  string    `db:"clerk_user_id"`
	FullName     string    `db:"full_name"`
	AvatarURL    string    `db:"avatar_url"`
	CreatedAt    time.Time `db:"created_at"`
}

type Item struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`

	Type        string    `db:"type"` // lost | found
	Title       string    `db:"title"`
	Description *string   `db:"description"`

	ImageURLs   []string  `db:"image_urls"`

	Location    string    `db:"location"`
	Campus      string    `db:"campus"`

	LostAt      *time.Time `db:"lost_at"`

	Tags        []string   `db:"tags"`

	IsConfirmed bool       `db:"is_confirmed"`

	CreatedAt   time.Time  `db:"created_at"`

	// Joined fields
	User        *User      `db:"-"`
}
