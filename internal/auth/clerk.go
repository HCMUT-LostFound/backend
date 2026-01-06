package auth

import (
	"fmt"
	"encoding/json"
	"os"
	"net/http"
	"strings"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

type ClerkVerifier struct {
	jwks *keyfunc.JWKS
	iss  string
}

func NewClerkVerifier(jwksURL, issuer string) (*ClerkVerifier, error) {
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			fmt.Println("JWKS refresh error:", err)
		},
	})
	if err != nil {
		return nil, err
	}

	return &ClerkVerifier{
		jwks: jwks,
		iss:  issuer,
	}, nil

}

func (v *ClerkVerifier) Verify(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return v.jwks.Keyfunc(token)
	})
	if err != nil {
		// ⛔ ignore "Token used before issued"
		if strings.Contains(err.Error(), "before issued") {
			// continue
		} else {
			return nil, nil, err
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("invalid claims")
	}

	// check issuer
	iss, ok := claims["iss"].(string)
	if !ok || iss != v.iss {
		return nil, nil, fmt.Errorf("invalid issuer")
	}

	return token, claims, nil
}



func ExtractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}
	return ""
}

func GetClerkUser(userID string) (*ClerkUser, error) {
	req, _ := http.NewRequest("GET", "https://api.clerk.com/v1/users/"+userID, nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("CLERK_SECRET_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("clerk api error: %s", resp.Status)
	}

	var user ClerkUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}