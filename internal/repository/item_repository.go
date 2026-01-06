package repository

import (
	"context"
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
	var items []Item

	query := `
	SELECT *
	FROM items
	WHERE is_confirmed = FALSE
	ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &items, query)
	return items, err
}

func (r *ItemRepository) ListByUser(ctx context.Context, clerkUserID string) ([]Item, error) {
	var items []Item

	query := `
	SELECT *
	FROM items
	WHERE user_id = $1
	  AND is_confirmed = FALSE
	ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &items, query, clerkUserID)
	return items, err
}

func (r *ItemRepository) Confirm(ctx context.Context, itemID uuid.UUID, clerkUserID string) error {
	query := `
	UPDATE items
	SET is_confirmed = TRUE
	WHERE id = $1
	  AND user_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, itemID, clerkUserID)
	return err
}
