package repository

import "github.com/jmoiron/sqlx"

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByClerkID(clerkID string) (*User, error) {
	var user User
	err := r.db.Get(&user,
		`SELECT * FROM users WHERE clerk_user_id = $1`,
		clerkID,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *User) error {
	query := `
	INSERT INTO users (clerk_user_id, full_name, avatar_url)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`
	return r.db.QueryRowx(
		query,
		user.ClerkUserID,
		user.FullName,
		user.AvatarURL,
	).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) UpsertProfile(
	clerkID, fullName, avatarURL string,
) error {
	query := `
	UPDATE users
	SET full_name = $1,
	    avatar_url = $2
	WHERE clerk_user_id = $3
	`
	_, err := r.db.Exec(query, fullName, avatarURL, clerkID)
	return err
}
