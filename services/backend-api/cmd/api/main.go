package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"noosphere/backend-api/internal/config"
	"noosphere/backend-api/internal/database"
	"noosphere/backend-api/internal/modules/auth"
	"noosphere/backend-api/internal/modules/chat"
	"noosphere/backend-api/internal/modules/user"
	"time"

	_ "noosphere/backend-api/docs" 
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2" 
)

// @title          Noosphere Backend API Engine
// @version        1.0
// @description    Core high-performance orchestration gateway for the Noosphere ecosystem.
// @host           localhost:8080
// @BasePath       /
func main() {
	cfg := config.LoadConfig()

	fmt.Println("Connecting to the database engine...")
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("CRITICAL: Database initialization failed: %v", err))
	}
	defer db.Close()
	fmt.Println("Database connection pool established successfully!")

	fmt.Println("Scanning for pending schema updates...")
	migrationsPath := "db/migrations" 
	err = database.RunMigrations(cfg.DatabaseURL, migrationsPath)
	if err != nil {
		panic(fmt.Sprintf("CRITICAL: Schema migration execution failed: %v", err))
	}
	fmt.Println("Database schemas are fully up-to-date!")

	chatRepo := chat.NewPostgresChatRepository(db)

	pythonAIEndpoint := "http://localhost:8000"
	aiClient := chat.NewHTTPInferenceClient(pythonAIEndpoint)
	
	chatService := chat.NewDefaultChatService(chatRepo, aiClient)
	chatController := chat.NewChatController(chatService) 

	// JWT Configuration & Fallback Generation
	var jwtSecret []byte
	if cfg.JWTSecret == "" {
		fmt.Println("WARNING: JWT_SECRET environment configuration is empty. Generating secure ephemeral key...")
		ephemeralSecret := make([]byte, 32)
		if _, err := rand.Read(ephemeralSecret); err != nil {
			panic(fmt.Sprintf("CRITICAL: Failed to generate secure random JWT ephemeral key: %v", err))
		}
		jwtSecret = ephemeralSecret
	} else {
		jwtSecret = []byte(cfg.JWTSecret)
	}

	accessExpiry, err := time.ParseDuration(cfg.JWTAccessExpiry)
	if err != nil {
		fmt.Printf("WARNING: Invalid JWT_ACCESS_EXPIRY value '%s', falling back to 15m: %v\n", cfg.JWTAccessExpiry, err)
		accessExpiry = 15 * time.Minute
	}

	refreshExpiry, err := time.ParseDuration(cfg.JWTRefreshExpiry)
	if err != nil {
		fmt.Printf("WARNING: Invalid JWT_REFRESH_EXPIRY value '%s', falling back to 168h: %v\n", cfg.JWTRefreshExpiry, err)
		refreshExpiry = 168 * time.Hour
	}

	userRepo := user.NewPostgresUserRepository(db)
	authRepo := auth.NewPostgresAuthRepository(db)
	authService := auth.NewDefaultAuthService(authRepo, userRepo, jwtSecret, accessExpiry, refreshExpiry)
	authController := auth.NewAuthController(authService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, 
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "noosphere-backend"}`))
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Post("/api/v1/chat/message", chatController.HandleSendMessage)
	r.Get("/api/v1/chat/session/{sessionID}/history", chatController.HandleGetChatHistory)

	r.Post("/api/v1/auth/register", authController.HandleRegister)
	r.Post("/api/v1/auth/login", authController.HandleLogin)
	r.Post("/api/v1/auth/refresh", authController.HandleRefresh)
	r.Post("/api/v1/auth/logout", authController.HandleLogout)

	serverAddress := ":" + cfg.Port
	fmt.Printf("Noosphere Engine online. Listening on port %s...\n", cfg.Port)
	
	// 3. Execution Engine Boot Loop
	err = http.ListenAndServe(serverAddress, r)
	if err != nil {
		panic(fmt.Sprintf("CRITICAL: Failed to start HTTP server: %v", err))
	}
}