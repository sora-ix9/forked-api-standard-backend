package httpclient

import (
	"time"

	"github.com/go-resty/resty/v2"
)

// NewClient initializes and returns a new Resty HTTP client
// Configured with sensible defaults for timeouts to prevent hanging connections
func NewClient() *resty.Client {
	client := resty.New()

	// Set default timeout
	client.SetTimeout(10 * time.Second)

	// Set retry count and wait time
	client.SetRetryCount(3)
	client.SetRetryWaitTime(100 * time.Millisecond)
	client.SetRetryMaxWaitTime(2 * time.Second)

	return client
}
