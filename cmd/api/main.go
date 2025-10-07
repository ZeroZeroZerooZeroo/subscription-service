package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/config"
	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/handler"
	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/repository"
	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/service"
	"github.com/ZeroZeroZerooZeroo/subscription-service/pkg/database"
)

func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		log.Printf("Warning: .env file not found")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}
}
func main() {

	loadEnvFile()
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
