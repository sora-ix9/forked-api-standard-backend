package stripe

import (
	"fdlp-standard-api/pkg/config"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
)

// StripeClient wraps the stripe-go client
type StripeClient struct {
	Client *client.API
	Config *config.Config
}

// NewClient initializes and returns a new Stripe client
func NewClient(cfg *config.Config) *StripeClient {
	stripe.Key = cfg.StripeSecretKey

	sc := &client.API{}
	sc.Init(cfg.StripeSecretKey, nil)

	// You can initialize specific service clients here if you want to wrap them
	// e.g., sc.PaymentIntents = paymentintent.Client{B: sc.Backend, Key: stripe.Key}

	return &StripeClient{
		Client: sc,
		Config: cfg,
	}
}
