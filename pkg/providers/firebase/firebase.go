package firebase

import (
	"context"
	"fdlp-standard-api/pkg/config"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// NewApp initializes and returns a new Firebase App
func NewApp(cfg *config.Config) *firebase.App {
	// If credentials file is not provided, we might look for default credentials
	// or maybe the user wants to configure it differently.
	// For now, we follow the pattern of using the config.

	var app *firebase.App
	var err error

	if cfg.FirebaseCredentialsFile != "" {
		opt := option.WithCredentialsFile(cfg.FirebaseCredentialsFile)
		app, err = firebase.NewApp(context.Background(), nil, opt)
	} else {
		// Try to initialize without File (Application Default Credentials)
		app, err = firebase.NewApp(context.Background(), nil)
	}

	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	return app
}
