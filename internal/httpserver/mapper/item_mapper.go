package mapper

import (
	"strings"
	"github.com/HCMUT-LostFound/backend/internal/auth"
	"github.com/HCMUT-LostFound/backend/internal/httpserver/dto"
	"github.com/HCMUT-LostFound/backend/internal/repository"
)

func ToItemResponse(item repository.Item) dto.ItemResponse {
	return dto.ItemResponse{
		ID:          item.ID,
		UserID:      item.UserID, // string → string
		Type:        item.Type,
		Title:       item.Title,
		Description: item.Description,
		ImageURLs:   item.ImageURLs,
		Location:    item.Location,
		Campus:      item.Campus,
		LostAt:      item.LostAt,
		Tags:        item.Tags,
		IsConfirmed: item.IsConfirmed,
		CreatedAt:   item.CreatedAt,
	}
}


func ToItemResponseWithReporter(item *repository.Item, user *auth.ClerkUser) dto.ItemResponse {
	fullName := strings.TrimSpace(user.FirstName + " " + user.LastName)

	return dto.ItemResponse{
		ID: item.ID,
		UserID: item.UserID, // string
		Title: item.Title,
		Type: item.Type,
		Description: item.Description,
		ImageURLs: item.ImageURLs,
		Location: item.Location,
		Campus: item.Campus,
		LostAt: item.LostAt,
		IsConfirmed: item.IsConfirmed,
		Reporter: &dto.Reporter{
			ID: item.UserID, // string
			FullName: fullName,
			ImageURL: user.ImageURL,
		},
	}
}


