package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChatMessageRepository struct {
	db *sqlx.DB
}

func NewChatMessageRepository(db *sqlx.DB) *ChatMessageRepository {
	return &ChatMessageRepository{db: db}
}

// Create tạo message mới
func (r *ChatMessageRepository) Create(ctx context.Context, message *ChatMessage) error {
	query := `
	INSERT INTO chat_messages (chat_id, sender_id, content)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`
	return r.db.QueryRowxContext(ctx, query,
		message.ChatID,
		message.SenderID,
		message.Content,
	).Scan(&message.ID, &message.CreatedAt)
}

// ListByChatID lấy danh sách messages của một chat
func (r *ChatMessageRepository) ListByChatID(ctx context.Context, chatID uuid.UUID) ([]ChatMessage, error) {
	var messages []ChatMessage
	
	query := `
	SELECT 
		cm.id, cm.chat_id, cm.sender_id, cm.content, cm.created_at,
		u.id as sender_user_id,
		u.full_name as sender_full_name,
		u.avatar_url as sender_avatar_url
	FROM chat_messages cm
	INNER JOIN users u ON cm.sender_id = u.id
	WHERE cm.chat_id = $1
	ORDER BY cm.created_at ASC
	`
	
	rows, err := r.db.QueryxContext(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var msg ChatMessage
		var sender User
		
		err := rows.Scan(
			&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.CreatedAt,
			&sender.ID, &sender.FullName, &sender.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
		
		msg.Sender = &sender
		messages = append(messages, msg)
	}
	
	return messages, nil
}

// GetLatestByChatID lấy message mới nhất của chat (dùng cho polling)
func (r *ChatMessageRepository) GetLatestByChatID(ctx context.Context, chatID uuid.UUID, afterMessageID *uuid.UUID) ([]ChatMessage, error) {
	var messages []ChatMessage
	
	query := `
	SELECT 
		cm.id, cm.chat_id, cm.sender_id, cm.content, cm.created_at,
		u.id as sender_user_id,
		u.full_name as sender_full_name,
		u.avatar_url as sender_avatar_url
	FROM chat_messages cm
	INNER JOIN users u ON cm.sender_id = u.id
	WHERE cm.chat_id = $1
	`
	args := []interface{}{chatID}
	
	if afterMessageID != nil {
		query += ` AND cm.id > $2`
		args = append(args, *afterMessageID)
	}
	
	query += ` ORDER BY cm.created_at ASC`
	
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var msg ChatMessage
		var sender User
		
		err := rows.Scan(
			&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.CreatedAt,
			&sender.ID, &sender.FullName, &sender.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
		
		msg.Sender = &sender
		messages = append(messages, msg)
	}
	
	return messages, nil
}

