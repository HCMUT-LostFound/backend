package mapper

import (
	"github.com/HCMUT-LostFound/backend/internal/httpserver/dto"
	"github.com/HCMUT-LostFound/backend/internal/repository"
)

func ToItemResponse(item repository.Item) dto.ItemResponse {
	res := dto.ItemResponse{
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

	// Map user if available
	if item.User != nil {
		res.User = &dto.UserResponse{
			ID:        item.User.ID,
			FullName:  item.User.FullName,
			AvatarURL: item.User.AvatarURL,
		}
	}

	return res
}
