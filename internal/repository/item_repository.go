package repository

import (
	"context"
	"database/sql"
	"time"
	"github.com/lib/pq"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ItemRepository struct {
	db *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Create(ctx context.Context, item *Item) error {
	query := `
	INSERT INTO items (
		user_id, type, title, description,
		image_urls, location, campus,
		lost_at, tags
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	RETURNING id, created_at
	`

	row := r.db.QueryRowContext(
		ctx,
		query,
		item.UserID,
		item.Type,
		item.Title,
		item.Description,
		pq.Array(item.ImageURLs), // 🔑 FIX
		item.Location,
		item.Campus,
		item.LostAt,
		pq.Array(item.Tags), // 🔑 FIX
	)

	return row.Scan(&item.ID, &item.CreatedAt)
}


func (r *ItemRepository) ListPublic(ctx context.Context) ([]Item, error) {
	type ItemRow struct {
		ID          uuid.UUID      `db:"id"`
		UserID      uuid.UUID      `db:"user_id"`
		Type        string         `db:"type"`
		Title       string         `db:"title"`
		Description sql.NullString `db:"description"`
		ImageURLs   pq.StringArray `db:"image_urls"`
		Location    string         `db:"location"`
		Campus      string         `db:"campus"`
		LostAt      sql.NullTime   `db:"lost_at"`
		Tags        pq.StringArray `db:"tags"`
		IsConfirmed bool           `db:"is_confirmed"`
		CreatedAt   time.Time      `db:"created_at"`
		// User fields
		UserIDFromUser uuid.UUID `db:"user_id_from_user"`
		UserFullName   string    `db:"user_full_name"`
		UserAvatarURL  string    `db:"user_avatar_url"`
	}

	var rows []ItemRow

	query := `
	SELECT 
		i.id, i.user_id, i.type, i.title, i.description,
		i.image_urls, i.location, i.campus, i.lost_at,
		i.tags, i.is_confirmed, i.created_at,
		u.id as user_id_from_user,
		u.full_name as user_full_name,
		u.avatar_url as user_avatar_url
	FROM items i
	INNER JOIN users u ON i.user_id = u.id
	WHERE i.is_confirmed = FALSE
	ORDER BY i.created_at DESC
	`

	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, err
	}

	// Convert to Item with User
	items := make([]Item, len(rows))
	for i, row := range rows {
		items[i] = Item{
			ID:          row.ID,
			UserID:      row.UserID,
			Type:        row.Type,
			Title:       row.Title,
			ImageURLs:   row.ImageURLs,
			Location:    row.Location,
			Campus:      row.Campus,
			Tags:        row.Tags,
			IsConfirmed: row.IsConfirmed,
			CreatedAt:   row.CreatedAt,
		}
		
		if row.Description.Valid {
			items[i].Description = &row.Description.String
		}
		if row.LostAt.Valid {
			items[i].LostAt = &row.LostAt.Time
		}
		
		// Add user info
		if row.UserFullName != "" {
			items[i].User = &User{
				ID:        row.UserIDFromUser,
				FullName:  row.UserFullName,
				AvatarURL: row.UserAvatarURL,
			}
		}
	}

	return items, nil
}

func (r *ItemRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]Item, error) {
	var items []Item

	query := `
	SELECT *
	FROM items
	WHERE user_id = $1
	  AND is_confirmed = FALSE
	ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &items, query, userID)
	return items, err
}

func (r *ItemRepository) Confirm(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) error {
	query := `
	UPDATE items
	SET is_confirmed = TRUE
	WHERE id = $1
	  AND user_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, itemID, userID)
	return err
}
