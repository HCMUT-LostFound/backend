package dto

import (
	"time"
	"github.com/google/uuid"
)

type ChatResponse struct {
	ID        uuid.UUID `json:"id"`
	ItemID    uuid.UUID `json:"itemId"`
	ItemTitle *string   `json:"itemTitle,omitempty"`
	ItemImage *string   `json:"itemImage,omitempty"`
	
	OtherUser *UserResponse `json:"otherUser,omitempty"`
	
	LastMessage *ChatMessageResponse `json:"lastMessage,omitempty"`
	
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ChatMessageResponse struct {
	ID        uuid.UUID    `json:"id"`
	ChatID    uuid.UUID    `json:"chatId"`
	SenderID  uuid.UUID    `json:"senderId"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"createdAt"`
	Sender    *UserResponse `json:"sender,omitempty"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"fullName"`
	AvatarURL string    `json:"avatarUrl"`
}

type CreateChatRequest struct {
	ItemID uuid.UUID `json:"itemId" binding:"required"`
}

type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

