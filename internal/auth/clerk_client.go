package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ClerkUser struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ImageURL  string `json:"image_url"`
	PrimaryEmailAddressID  string `json:"primary_email_address_id"`
	EmailAddresses []struct {
		ID           string `json:"id"`
		EmailAddress string `json:"email_address"`
	} `json:"email_addresses"`
}

func (u *ClerkUser) PrimaryEmail() string {
	for _, e := range u.EmailAddresses {
		if e.ID == u.PrimaryEmailAddressID {
			return e.EmailAddress
		}
	}
	return ""
}

func FetchClerkUser(clerkUserID, secretKey string) (*ClerkUser, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.clerk.com/v1/users/%s", clerkUserID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("clerk api error: %s", resp.Status)
	}

	var user ClerkUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
