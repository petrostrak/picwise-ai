package handler

import (
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
		SuccessURL: stripe.String(""),
		CancelURL:  stripe.String(""),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("prod_id"),
				Quantity: stripe.Int64(1),
			},
		},
	}

	s, err := session.New(checkoutParams)
	if err != nil {
		return err
	}

	http.Redirect(w, r, s.URL, http.StatusSeeOther)
	return nil
}
