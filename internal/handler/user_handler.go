package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HCMUT-LostFound/backend/internal/repository"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GET /api/me
func (h *UserHandler) GetMe(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found in context",
		})
		return
	}

	user, ok := userValue.(*repository.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user type",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
