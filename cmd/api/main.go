package main

import (
	"log"
	"net/http"

	_ "github.com/ZeroZeroZerooZeroo/subscription-service/docs"
	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/config"
	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/handler"
	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/repository"
	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/service"
	"github.com/ZeroZeroZerooZeroo/subscription-service/pkg/database"
)

// @title Subscription Service API
// @version 1.0
// @description Сервис управления подписками с расчетом стоимости
// @host localhost:8080
func main() {

	log.Println("Starting subscription service")
	cfg := config.LoadConfig()

	db, err := database.NewPostgres(cfg.Database.GetDBConnectionString())

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer func() {
		log.Println("Closing database connection")
		db.Close()
	}()

	if err := database.RunMigrations(cfg.Database.GetMigrationConnectionString()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	repo := repository.NewSubscriptionRepository(db.DB)
	svc := service.NewSubscriptionService(repo)
	handler := handler.NewSubscriptionHandler(svc)

	mux := http.NewServeMux()
	handler.SetupRoutes(mux)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("Server stopped")
}
