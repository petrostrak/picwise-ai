package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/petrostrak/picwise-ai/view/credits"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"net/http"
	"os"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, credits.Index())
}

func HandleStripeCheckoutCreate(w http.ResponseWriter, r *http.Request) error {
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	checkoutParams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("http://localhost:3000/checkout/success"),
		CancelURL:  stripe.String("http://localhost:3000/checkout/cancel"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(chi.URLParam(r, "priceID")),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}

	s, err := session.New(checkoutParams)
	if err != nil {
		return err
	}

	return hxRedirect(w, r, s.URL)
}

func HandleStripeCheckoutSuccess(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func HandleStripeCheckoutCancel(w http.ResponseWriter, r *http.Request) error {
	return nil
}
