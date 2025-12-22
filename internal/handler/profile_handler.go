package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HCMUT-LostFound/backend/internal/repository"
)

type ProfileHandler struct{}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found in context",
		})
		return
	}

	user := userValue.(*repository.User)

	c.JSON(http.StatusOK, gin.H{
		"fullName":  user.FullName,
		"avatarUrl": user.AvatarURL,
	})
}
