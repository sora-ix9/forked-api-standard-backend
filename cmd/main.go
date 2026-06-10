package main

import (
	"fdlp-standard-api/internal/echo"
	"fdlp-standard-api/internal/utils"
	"fdlp-standard-api/pkg/config"
	"fdlp-standard-api/pkg/db"
	"fdlp-standard-api/pkg/redisclient"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load("env_files/.env.dev")
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Global settings
	utils.SetGlobalTimezone()
	log.Println("Global timezone set to Asia/Bangkok (UTC+07:00)")

	// Initialize configuration
	cfg := config.New()

	// Initialize Databases
	mongodb := db.NewMongoDB(10)
	defer mongodb.Disconnect()

	redisClient := redisclient.NewClient(cfg)
	defer redisClient.Close()

	// Start Echo Server
	echo.StartServer(cfg, mongodb, redisClient)
}
