package middleware

import (
	"net/http"
	"strings"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/HCMUT-LostFound/backend/internal/auth"
)

var allowedDomains = []string{
	"hcmut.edu.vn",
}

func DomainGuard() gin.HandlerFunc {
	return func(c *gin.Context) {

		clerkIDAny, exists := c.Get("clerk_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no clerk user id"})
			return
		}

		clerkID := clerkIDAny.(string)

		clerkUser, err := auth.FetchClerkUser(clerkID, os.Getenv("CLERK_SECRET_KEY"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "cannot fetch clerk user"})
			return
		}

		email := clerkUser.PrimaryEmail()
		if email == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no email"})
			return
		}

		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid email"})
			return
		}

		domain := parts[1]
		for _, d := range allowedDomains {
			if domain == d {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":  "email domain not allowed",
			"email":  email,
			"domain": domain,
		})
	}
}
