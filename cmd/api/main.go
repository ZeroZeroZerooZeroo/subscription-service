package main

import (
	"log"

	"github.com/ZeroZeroZerooZeroo/subscription-service/internal/config"
	"github.com/ZeroZeroZerooZeroo/subscription-service/pkg/database"
)

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

}
