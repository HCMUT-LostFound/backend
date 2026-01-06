package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/HCMUT-LostFound/backend/internal/httpserver/dto"
	"github.com/HCMUT-LostFound/backend/internal/httpserver/mapper"
	"github.com/HCMUT-LostFound/backend/internal/repository"
)

type ChatHandler struct {
	chatRepo    *repository.ChatRepository
	messageRepo *repository.ChatMessageRepository
}

func NewChatHandler(
	chatRepo *repository.ChatRepository,
	messageRepo *repository.ChatMessageRepository,
) *ChatHandler {
	return &ChatHandler{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
	}
}

// POST /api/chats - Tạo hoặc lấy chat về một item
func (h *ChatHandler) CreateOrGet(c *gin.Context) {
	var req dto.CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(*repository.User)

	chat, err := h.chatRepo.CreateOrGet(c.Request.Context(), req.ItemID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapper.ToChatResponse(*chat))
}

// GET /api/chats - Lấy danh sách chats của user
func (h *ChatHandler) List(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)

	chats, err := h.chatRepo.ListByUser(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := make([]dto.ChatResponse, 0, len(chats))
	for _, chat := range chats {
		res = append(res, mapper.ToChatResponse(chat))
	}

	c.JSON(http.StatusOK, res)
}

// GET /api/chats/:id/messages - Lấy messages của một chat
func (h *ChatHandler) GetMessages(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)

	chatIDParam := c.Param("id")
	chatID, err := uuid.Parse(chatIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
		return
	}

	// Kiểm tra user có quyền truy cập chat
	_, err = h.chatRepo.GetByID(c.Request.Context(), chatID, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "chat not found"})
		return
	}

	// Lấy after parameter (optional, cho polling)
	afterParam := c.Query("after")
	var afterMessageID *uuid.UUID
	if afterParam != "" {
		afterID, err := uuid.Parse(afterParam)
		if err == nil {
			afterMessageID = &afterID
		}
	}

	var messages []repository.ChatMessage
	if afterMessageID != nil {
		messages, err = h.messageRepo.GetLatestByChatID(c.Request.Context(), chatID, afterMessageID)
	} else {
		messages, err = h.messageRepo.ListByChatID(c.Request.Context(), chatID)
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := make([]dto.ChatMessageResponse, 0, len(messages))
	for _, msg := range messages {
		res = append(res, mapper.ToChatMessageResponse(msg))
	}

	c.JSON(http.StatusOK, res)
}

// POST /api/chats/:id/messages - Gửi message
func (h *ChatHandler) SendMessage(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)

	chatIDParam := c.Param("id")
	chatID, err := uuid.Parse(chatIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
		return
	}

	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kiểm tra user có quyền truy cập chat
	_, err = h.chatRepo.GetByID(c.Request.Context(), chatID, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "chat not found"})
		return
	}

	message := &repository.ChatMessage{
		ChatID:   chatID,
		SenderID: user.ID,
		Content:  req.Content,
	}

	if err := h.messageRepo.Create(c.Request.Context(), message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Load lại với sender info
	messages, err := h.messageRepo.ListByChatID(c.Request.Context(), chatID)
	if err == nil && len(messages) > 0 {
		for _, m := range messages {
			if m.ID == message.ID {
				message = &m
				break
			}
		}
	}

	c.JSON(http.StatusCreated, mapper.ToChatMessageResponse(*message))
}

