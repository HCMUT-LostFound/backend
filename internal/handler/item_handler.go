package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"

	"github.com/HCMUT-LostFound/backend/internal/httpserver/dto"
	"github.com/HCMUT-LostFound/backend/internal/httpserver/mapper"
	"github.com/HCMUT-LostFound/backend/internal/repository"
)

type ItemHandler struct {
	repo *repository.ItemRepository
}

func NewItemHandler(repo *repository.ItemRepository) *ItemHandler {
	return &ItemHandler{repo: repo}
}

func validateCreateItem(req *dto.CreateItemRequest) error {
	if req.Type != "lost" && req.Type != "found" {
		return gin.Error{
			Err:  http.ErrNotSupported,
			Type: gin.ErrorTypeBind,
		}
	}

	if req.Type == "lost" && len(req.ImageURLs) == 0 {
		return gin.Error{
			Err:  gin.Error{Err: http.ErrMissingFile},
			Type: gin.ErrorTypeBind,
		}
	}

	return nil
}

func (h *ItemHandler) Create(c *gin.Context) {
	var req dto.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateCreateItem(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item data"})
		return
	}

	user := c.MustGet("user").(*repository.User)

	item := &repository.Item{
		UserID:    user.ID,
		Type:      req.Type,
		Title:     req.Title,
		Description: req.Description,
		ImageURLs: req.ImageURLs,
		Location:  req.Location,
		Campus:    req.Campus,
		LostAt:    req.LostAt,
		Tags:      req.Tags,
	}

	if err := h.repo.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) ListPublic(c *gin.Context) {
	items, err := h.repo.ListPublic(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := make([]dto.ItemResponse, 0, len(items))
	for _, item := range items {
		res = append(res, mapper.ToItemResponse(item))
	}

	c.JSON(http.StatusOK, res)
}
