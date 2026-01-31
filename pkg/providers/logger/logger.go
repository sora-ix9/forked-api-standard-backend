package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps the logrus logger
type Logger struct {
	*logrus.Logger
}

// NewLogger initializes and returns a new Logger
// You might want to pass config here to determine format (JSON vs Text) and Level
func NewLogger() *Logger {
	l := logrus.New()

	// Default to JSON format for production-readiness
	// In a real scenario, check environment variable (e.g. APP_ENV)
	// If DEV -> TextFormatter, if PROD -> JSONFormatter
	env := os.Getenv("APP_ENV")
	if env == "production" {
		l.SetFormatter(&logrus.JSONFormatter{})
	} else {
		l.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.InfoLevel)

	return &Logger{l}
}
