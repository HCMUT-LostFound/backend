package auth

import (
	"fmt"
	"net/http"

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
	token, err := jwt.Parse(tokenString, v.jwks.Keyfunc)
	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}

	// check issuer manually (Clerk)
	if claims["iss"] != v.iss {
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
