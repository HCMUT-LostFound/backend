package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/HCMUT-LostFound/backend/internal/handler"
)
type Dependencies struct {
	UserHandler *handler.UserHandler
	ProfileHandler *handler.ProfileHandler
	ItemHandler *handler.ItemHandler
}
func RegisterRoutes(
	r *gin.Engine,
	public *gin.RouterGroup,
	protected *gin.RouterGroup,
	deps *Dependencies,
) {
	// public
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	public.GET("/items", deps.ItemHandler.ListPublic)
	// protected
	protected.GET("/me", deps.UserHandler.GetMe)
	protected.GET("/profile", deps.ProfileHandler.GetProfile)
	protected.POST("/items", deps.ItemHandler.Create)
	// protected.GET("/items", deps.ItemHandler.ListPublic)
}
