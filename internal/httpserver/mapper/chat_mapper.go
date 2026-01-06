package mapper

import (
	"github.com/HCMUT-LostFound/backend/internal/httpserver/dto"
	"github.com/HCMUT-LostFound/backend/internal/repository"
)

func ToChatResponse(chat repository.Chat) dto.ChatResponse {
	res := dto.ChatResponse{
		ID:        chat.ID,
		ItemID:    chat.ItemID,
		ItemTitle: chat.ItemTitle,
		ItemImage: chat.ItemImage,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}
	
	if chat.OtherUser != nil {
		res.OtherUser = &dto.UserResponse{
			ID:        chat.OtherUser.ID,
			FullName:  chat.OtherUser.FullName,
			AvatarURL: chat.OtherUser.AvatarURL,
		}
	}
	
	if chat.LastMessage != nil {
		res.LastMessage = &dto.ChatMessageResponse{
			ID:        chat.LastMessage.ID,
			Content:   chat.LastMessage.Content,
			CreatedAt: chat.LastMessage.CreatedAt,
		}
	}
	
	return res
}

func ToChatMessageResponse(msg repository.ChatMessage) dto.ChatMessageResponse {
	res := dto.ChatMessageResponse{
		ID:        msg.ID,
		ChatID:    msg.ChatID,
		SenderID:  msg.SenderID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
	}
	
	if msg.Sender != nil {
		res.Sender = &dto.UserResponse{
			ID:        msg.Sender.ID,
			FullName:  msg.Sender.FullName,
			AvatarURL: msg.Sender.AvatarURL,
		}
	}
	
	return res
}

