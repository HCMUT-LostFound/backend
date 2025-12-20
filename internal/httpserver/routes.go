package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/HCMUT-LostFound/backend/internal/handler"
)
type Dependencies struct {
	UserHandler *handler.UserHandler
}
func RegisterRoutes(
	r *gin.Engine,
	protected *gin.RouterGroup,
	deps *Dependencies,
) {
	// public
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// protected
	protected.GET("/me", deps.UserHandler.GetMe)
}
