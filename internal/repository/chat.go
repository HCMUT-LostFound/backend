package repository

import (
	"time"
	"github.com/google/uuid"
)

type Chat struct {
	ID          uuid.UUID `db:"id"`
	ItemID      uuid.UUID `db:"item_id"`
	InitiatorID uuid.UUID `db:"initiator_id"`
	ItemOwnerID uuid.UUID `db:"item_owner_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	
	// Joined fields (optional, for responses)
	ItemTitle   *string `db:"item_title"`
	ItemImage   *string `db:"item_image"`
	OtherUser   *User   `db:"other_user"`
	LastMessage *ChatMessage `db:"last_message"`
}

type ChatMessage struct {
	ID        uuid.UUID `db:"id"`
	ChatID    uuid.UUID `db:"chat_id"`
	SenderID  uuid.UUID `db:"sender_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	
	// Joined fields (optional)
	Sender *User `db:"sender"`
}

