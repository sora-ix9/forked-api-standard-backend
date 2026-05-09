package config

import (
	"log"
	"os"
)

type Config struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDb       string

	JWTSecret string

	// S3 Config
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSBucketName      string

	// Firebase Config
	FirebaseCredentialsFile string

	// Email Config
	SMTPHost      string
	SMTPPort      string
	SMTPUsername  string
	SMTPPassword  string
	SMTPFromEmail string

	// Stripe Config
	StripeSecretKey      string
	StripePublishableKey string
	StripeWebhookSecret  string

	// MongoDB Config
	MongoDBURI string
	MongoDBDb  string
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", "redispassword"),
		RedisDb:       getEnv("REDIS_DB", "0"),
		JWTSecret:     getEnv("JWT_SECRET", "ABCD"),

		// S3 Config
		AWSRegion:          getEnv("AWS_REGION", "ap-southeast-1"),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSBucketName:      getEnv("AWS_BUCKET_NAME", ""),

		// Firebase Config
		FirebaseCredentialsFile: getEnv("FIREBASE_CREDENTIALS_FILE", ""),

		// Email Config
		SMTPHost:      getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:      getEnv("SMTP_PORT", "587"),
		SMTPUsername:  getEnv("SMTP_USERNAME", "user@example.com"),
		SMTPPassword:  getEnv("SMTP_PASSWORD", "password"),
		SMTPFromEmail: getEnv("SMTP_FROM_EMAIL", "no-reply@example.com"),

		// Stripe Config
		StripeSecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
		StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		StripeWebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),

		// MongoDB Config
		MongoDBURI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBDb:  getEnv("MONGODB_DB", "fdlp_mongo"),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Warning: %s is not set. Using default value: %s", key, defaultValue)

	return defaultValue
}
