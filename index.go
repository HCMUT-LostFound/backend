package backend

// Export Handler for Vercel

import (
	"log"
	"net/http"
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
	apphandler "github.com/HCMUT-LostFound/backend/internal/handler"
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
		allowedOrigins = append(allowedOrigins, prodOrigins)
	}

	router.Use(cors.New(cors.Config{
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
		log.Printf("Warning: Failed to create Clerk verifier: %v", err)
		// Don't fatal in serverless - will be handled at request time
		return
	}

	userRepo := repository.NewUserRepository(database)
	itemRepo := repository.NewItemRepository(database)
	chatRepo := repository.NewChatRepository(database)
	chatMessageRepo := repository.NewChatMessageRepository(database)

	userHandler := apphandler.NewUserHandler()
	profileHandler := apphandler.NewProfileHandler()
	itemHandler := apphandler.NewItemHandler(itemRepo)
	chatHandler := apphandler.NewChatHandler(chatRepo, chatMessageRepo)

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

func init() {
	initRouter()
}

// Handler is the entry point for Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	if router == nil {
		initRouter()
	}
	router.ServeHTTP(w, r)
}

