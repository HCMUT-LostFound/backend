package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChatRepository struct {
	db *sqlx.DB
}

func NewChatRepository(db *sqlx.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// CreateOrGet tạo chat mới hoặc trả về chat đã tồn tại
func (r *ChatRepository) CreateOrGet(ctx context.Context, itemID, initiatorID uuid.UUID) (*Chat, error) {
	// Lấy item owner
	var itemOwnerID uuid.UUID
	err := r.db.GetContext(ctx, &itemOwnerID,
		`SELECT user_id FROM items WHERE id = $1`,
		itemID,
	)
	if err != nil {
		return nil, err
	}

	// Kiểm tra chat đã tồn tại chưa
	var chat Chat
	err = r.db.GetContext(ctx, &chat,
		`SELECT * FROM chats WHERE item_id = $1 AND initiator_id = $2`,
		itemID, initiatorID,
	)
	
	if err == sql.ErrNoRows {
		// Tạo chat mới
		chat = Chat{
			ItemID:      itemID,
			InitiatorID: initiatorID,
			ItemOwnerID: itemOwnerID,
		}
		err = r.db.QueryRowxContext(ctx,
			`INSERT INTO chats (item_id, initiator_id, item_owner_id)
			 VALUES ($1, $2, $3)
			 RETURNING id, created_at, updated_at`,
			chat.ItemID, chat.InitiatorID, chat.ItemOwnerID,
		).Scan(&chat.ID, &chat.CreatedAt, &chat.UpdatedAt)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	// Chat đã tồn tại, đã load vào biến chat từ GetContext ở trên
	
	return &chat, nil
}

// ListByUser lấy danh sách chats của user
func (r *ChatRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]Chat, error) {
	var chats []Chat
	
	query := `
	SELECT 
		c.id, c.item_id, c.initiator_id, c.item_owner_id, c.created_at, c.updated_at,
		i.title as item_title,
		COALESCE(i.image_urls[1], '') as item_image,
		CASE 
			WHEN c.initiator_id = $1 THEN u2.id
			ELSE u1.id
		END as other_user_id,
		CASE 
			WHEN c.initiator_id = $1 THEN u2.full_name
			ELSE u1.full_name
		END as other_user_full_name,
		CASE 
			WHEN c.initiator_id = $1 THEN COALESCE(u2.avatar_url, '')
			ELSE COALESCE(u1.avatar_url, '')
		END as other_user_avatar_url,
		cm.id as last_message_id,
		cm.content as last_message_content,
		cm.created_at as last_message_created_at
	FROM chats c
	INNER JOIN items i ON c.item_id = i.id
	INNER JOIN users u1 ON c.initiator_id = u1.id
	INNER JOIN users u2 ON c.item_owner_id = u2.id
	LEFT JOIN LATERAL (
		SELECT id, content, created_at
		FROM chat_messages
		WHERE chat_id = c.id
		ORDER BY created_at DESC
		LIMIT 1
	) cm ON true
	WHERE c.initiator_id = $1 OR c.item_owner_id = $1
	ORDER BY c.updated_at DESC
	`
	
	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var chat Chat
		var otherUserID uuid.UUID
		var otherUserName, otherUserAvatar sql.NullString
		var itemTitle, itemImage sql.NullString
		var lastMsgID sql.NullString
		var lastMsgContent sql.NullString
		var lastMsgTime sql.NullTime
		
		err := rows.Scan(
			&chat.ID, &chat.ItemID, &chat.InitiatorID, &chat.ItemOwnerID,
			&chat.CreatedAt, &chat.UpdatedAt,
			&itemTitle, &itemImage,
			&otherUserID, &otherUserName, &otherUserAvatar,
			&lastMsgID, &lastMsgContent, &lastMsgTime,
		)
		if err != nil {
			return nil, err
		}
		
		if itemTitle.Valid {
			chat.ItemTitle = &itemTitle.String
		}
		if itemImage.Valid {
			chat.ItemImage = &itemImage.String
		}
		
		// Set other user
		if otherUserName.Valid {
			chat.OtherUser = &User{
				ID:        otherUserID,
				FullName:  otherUserName.String,
				AvatarURL: otherUserAvatar.String,
			}
		}
		
		// Set last message
		if lastMsgID.Valid && lastMsgContent.Valid {
			msgID, _ := uuid.Parse(lastMsgID.String)
			chat.LastMessage = &ChatMessage{
				ID:        msgID,
				Content:   lastMsgContent.String,
				CreatedAt: lastMsgTime.Time,
			}
		}
		
		chats = append(chats, chat)
	}
	
	return chats, nil
}

// GetByID lấy chat theo ID (với kiểm tra quyền)
func (r *ChatRepository) GetByID(ctx context.Context, chatID, userID uuid.UUID) (*Chat, error) {
	var chat Chat
	var itemTitle, itemImage sql.NullString
	err := r.db.QueryRowxContext(ctx,
		`SELECT c.*, i.title as item_title, COALESCE(i.image_urls[1], '') as item_image
		 FROM chats c
		 INNER JOIN items i ON c.item_id = i.id
		 WHERE c.id = $1 AND (c.initiator_id = $2 OR c.item_owner_id = $2)`,
		chatID, userID,
	).Scan(
		&chat.ID, &chat.ItemID, &chat.InitiatorID, &chat.ItemOwnerID,
		&chat.CreatedAt, &chat.UpdatedAt,
		&itemTitle, &itemImage,
	)
	if err != nil {
		return nil, err
	}
	
	if itemTitle.Valid {
		chat.ItemTitle = &itemTitle.String
	}
	if itemImage.Valid {
		chat.ItemImage = &itemImage.String
	}
	
	return &chat, nil
}

