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
