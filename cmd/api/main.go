package main

import (
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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
	
	// CORS middleware - Read from environment or use defaults
	allowedOrigins := []string{
		"http://localhost:8081", "http://localhost:8082", "http://localhost:19006",
		"http://192.168.1.69:8081", "http://192.168.1.69:8082", "http://192.168.1.69:19006",
		"http://192.168.1.88:8081", "http://192.168.1.88:8082", "http://192.168.1.88:19006",
		"http://172.28.144.1:8081", "http://172.28.144.1:8082", "http://172.28.144.1:19006",
	}
	
	// Add production origins from environment if set (comma-separated)
	if prodOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); prodOrigins != "" {
		// In production, you can set: CORS_ALLOWED_ORIGINS=https://your-domain.com,exp://your-expo-url
		// For now, we'll keep development origins and add production ones
		// Note: In production, you might want to replace the entire list
		allowedOrigins = append(allowedOrigins, prodOrigins)
	}
	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	verifier, err := auth.NewClerkVerifier(
		os.Getenv("CLERK_JWKS_URL"),
		os.Getenv("CLERK_ISSUER"),
	)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(database)
	itemRepo := repository.NewItemRepository(database)
	chatRepo := repository.NewChatRepository(database)
	chatMessageRepo := repository.NewChatMessageRepository(database)

	userHandler := handler.NewUserHandler()
	profileHandler := handler.NewProfileHandler()
	itemHandler := handler.NewItemHandler(itemRepo)
	chatHandler := handler.NewChatHandler(chatRepo, chatMessageRepo)

	// ===== PUBLIC API GROUP =====
	public := r.Group("/api")

	// ===== PROTECTED API GROUP =====
	protected := r.Group("/api")
	protected.Use(middleware.ClerkAuth(verifier, userRepo))

	httpserver.RegisterRoutes(
		r,
		public,
		protected,
		&httpserver.Dependencies{
			UserHandler:    userHandler,
			ProfileHandler: profileHandler,
			ItemHandler:    itemHandler,
			ChatHandler:    chatHandler,
		},
	)

	log.Printf("API listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
