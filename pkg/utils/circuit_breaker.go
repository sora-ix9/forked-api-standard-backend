package utils

import (
	"time"

	"github.com/sony/gobreaker"
)

// NewCircuitBreaker creates a new circuit breaker with default settings
func NewCircuitBreaker(name string) *gobreaker.TwoStepCircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}

	return gobreaker.NewTwoStepCircuitBreaker(settings)
}

// ExecuteWithBreaker executes a function using a circuit breaker
func ExecuteWithBreaker(cb *gobreaker.TwoStepCircuitBreaker, fn func() (interface{}, error)) (interface{}, error) {
	success, err := cb.Allow()
	if err != nil {
		return nil, err
	}

	result, err := fn()
	success(err == nil)

	return result, err
}
