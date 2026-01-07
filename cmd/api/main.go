package main

import (
	"log"
	"os"
	"strings"
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

var router *gin.Engine

func initRouter() {
	_ = godotenv.Load()
	cfg := config.Load()
	database := db.NewPostgres(cfg.DBUrl)
	// Note: Don't close database in serverless - connections are managed by pool

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// CORS middleware
	allowedOrigins := []string{
		"http://localhost:8081", "http://localhost:8082", "http://localhost:19006",
		"http://192.168.1.69:8081", "http://192.168.1.69:8082", "http://192.168.1.69:19006",
		"http://192.168.1.88:8081", "http://192.168.1.88:8082", "http://192.168.1.88:19006",
		"http://172.28.144.1:8081", "http://172.28.144.1:8082", "http://172.28.144.1:19006",
	}

	if prodOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); prodOrigins != "" {
		// Split by comma if multiple origins
		origins := strings.Split(prodOrigins, ",")
		for _, origin := range origins {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				allowedOrigins = append(allowedOrigins, origin)
			}
		}
	}

	// For mobile apps (React Native/Expo), allow all origins
	// Mobile apps don't have a specific origin, so we need to allow all
	allowAllOrigins := os.Getenv("ALLOW_ALL_ORIGINS") == "true"
	
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	if allowAllOrigins {
		// Allow all origins for mobile apps
		corsConfig.AllowOriginFunc = func(origin string) bool {
			return true
		}
		log.Println("[CORS] Allowing all origins (mobile app mode)")
	} else {
		corsConfig.AllowOrigins = allowedOrigins
		log.Printf("[CORS] Allowing specific origins: %v", allowedOrigins)
	}

	router.Use(cors.New(corsConfig))

	verifier, err := auth.NewClerkVerifier(
		os.Getenv("CLERK_JWKS_URL"),
		os.Getenv("CLERK_ISSUER"),
	)
	if err != nil {
		log.Printf("Warning: Failed to create Clerk verifier: %v", err)
		// Don't fatal - will be handled at request time
		return
	}

	userRepo := repository.NewUserRepository(database)
	itemRepo := repository.NewItemRepository(database)
	chatRepo := repository.NewChatRepository(database)
	chatMessageRepo := repository.NewChatMessageRepository(database)

	userHandler := handler.NewUserHandler()
	profileHandler := handler.NewProfileHandler()
	itemHandler := handler.NewItemHandler(itemRepo)
	chatHandler := handler.NewChatHandler(chatRepo, chatMessageRepo)

	public := router.Group("/api")
	protected := router.Group("/api")
	protected.Use(middleware.ClerkAuth(verifier, userRepo))

	httpserver.RegisterRoutes(
		router,
		public,
		protected,
		&httpserver.Dependencies{
			UserHandler:    userHandler,
			ProfileHandler: profileHandler,
			ItemHandler:    itemHandler,
			ChatHandler:    chatHandler,
		},
	)
}

func main() {
	if router == nil {
		initRouter()
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
