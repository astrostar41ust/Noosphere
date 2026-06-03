package main

import (
	"fmt"
	"net/http"
	"noosphere/backend-api/internal/config"
	"noosphere/backend-api/internal/database"
	"noosphere/backend-api/internal/modules/chat"

	_ "noosphere/backend-api/docs" 
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "noosphere-backend"}`))
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Post("/api/v1/chat/message", chatController.HandleSendMessage)

	serverAddress := ":" + cfg.Port
	fmt.Printf("Noosphere Engine online. Listening on port %s...\n", cfg.Port)
	
	// 3. Execution Engine Boot Loop
	err = http.ListenAndServe(serverAddress, r)
	if err != nil {
		panic(fmt.Sprintf("CRITICAL: Failed to start HTTP server: %v", err))
	}
}