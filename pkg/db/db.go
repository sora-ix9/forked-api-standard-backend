package db

import (
	"fdlp-standard-api/pkg/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitializeDB initializes and returns a database connection
func InitializeDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
		cfg.DBTimeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Get the underlying sql.DB object to configure the connection pool
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database: " + err.Error())
	}

	// Set the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)

	// Set the maximum number of idle connections to the database
	sqlDB.SetMaxIdleConns(100)

	// Set the maximum lifetime of a connection to 5 minutes
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db
}

func CloseDB(db *gorm.DB) {
	// Get the underlying sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database: %v", err)
	}

	// Close the database connection
	err = sqlDB.Close()
	if err != nil {
		log.Fatalf("Failed to close database: %v", err)
	}

	log.Println("Database connection closed successfully")
}
