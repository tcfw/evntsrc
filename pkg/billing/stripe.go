package billing

import (
	"log"
	"os"

	stripe "github.com/stripe/stripe-go"
)

func init() {
	stripeKey, ok := os.LookupEnv("STRIPE_KEY")
	if !ok {
		panic("No stripe key provided")
	}
	stripe.Key = stripeKey
	log.Println("Stripe key set")
}
