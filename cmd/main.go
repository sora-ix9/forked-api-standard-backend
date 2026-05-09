package main

import (
	"fdlp-standard-api/internal/echo"
	"fdlp-standard-api/internal/utils"
	"fdlp-standard-api/pkg/config"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Loading .env from the root directory
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: Error loading .env file, using default environment variables")
	}

	utils.SetGlobalTimezone()
	log.Println("Global timezone set to Asia/Bangkok (UTC+07:00)")

	// Initialize Configuration
	cfg := config.New()

	// Start Server
	echo.InitServer(cfg)
}
