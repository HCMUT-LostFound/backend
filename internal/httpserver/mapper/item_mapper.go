package mapper

import (
	"github.com/HCMUT-LostFound/backend/internal/httpserver/dto"
	"github.com/HCMUT-LostFound/backend/internal/repository"
)

func ToItemResponse(item repository.Item) dto.ItemResponse {
	return dto.ItemResponse{
		ID:          item.ID,
		UserID:      item.UserID,
		Type:        item.Type,
		Title:       item.Title,
		Description: item.Description,
		ImageURLs:   item.ImageURLs,
		Location:    item.Location,
		Campus:      item.Campus,
		LostAt:      item.LostAt,
		Tags:        item.Tags,
		CreatedAt:   item.CreatedAt,
	}
}
