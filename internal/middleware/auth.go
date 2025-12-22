package middleware

import (
	"log"
	"os"
	"strings"
	"net/http"
	"github.com/gin-gonic/gin"

	"github.com/HCMUT-LostFound/backend/internal/auth"
	"github.com/HCMUT-LostFound/backend/internal/repository"
)

func ClerkAuth(
	verifier *auth.ClerkVerifier,
	userRepo *repository.UserRepository,
) gin.HandlerFunc {

	return func(c *gin.Context) {
		tokenString := auth.ExtractBearerToken(c.Request)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		_, claims, err := verifier.Verify(tokenString)
		if err != nil {
			log.Println("[AUTH] verify error:", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		clerkID := claims["sub"].(string)
		user, err := userRepo.GetByClerkID(clerkID)
		if err != nil {
			user = &repository.User{
				ClerkUserID: clerkID,
			}
			_ = userRepo.Create(user)
		}
		clerkUser, err := auth.FetchClerkUser(
			clerkID,
			os.Getenv("CLERK_SECRET_KEY"),
		)
		if err == nil {
			fullName := strings.TrimSpace(
				clerkUser.FirstName + " " + clerkUser.LastName,
			)

			_ = userRepo.UpsertProfile(
				clerkID,
				fullName,
				clerkUser.ImageURL,
			)
		}

		user, _ = userRepo.GetByClerkID(clerkID)

		c.Set("user", user)
		c.Next()
	}
}
