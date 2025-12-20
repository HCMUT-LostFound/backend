package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/HCMUT-LostFound/backend/internal/config"
	"github.com/HCMUT-LostFound/backend/internal/db"
	"github.com/HCMUT-LostFound/backend/internal/httpserver"
	"github.com/HCMUT-LostFound/backend/internal/auth"
	"github.com/HCMUT-LostFound/backend/internal/repository"
	"github.com/HCMUT-LostFound/backend/internal/middleware"
	"github.com/HCMUT-LostFound/backend/internal/handler"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	database := db.NewPostgres(cfg.DBUrl)
	defer database.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	verifier, err := auth.NewClerkVerifier(
	os.Getenv("CLERK_JWKS_URL"),
	os.Getenv("CLERK_ISSUER"),
	)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(database)

	protected := r.Group("/api")
	protected.Use(middleware.ClerkAuth(verifier, userRepo))
	userHandler := handler.NewUserHandler()

	httpserver.RegisterRoutes(
		r,
		protected,
		&httpserver.Dependencies{
			UserHandler: userHandler,
		},
	)

	log.Printf("API listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
