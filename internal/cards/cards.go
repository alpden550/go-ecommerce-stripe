package cards

import (
	"errors"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/paymentmethod"
	"github.com/stripe/stripe-go/v74/refund"
	"github.com/stripe/stripe-go/v74/subscription"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	TransactionStatusId int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

func (c *Card) GetPaymentMethod(id string) (*stripe.PaymentMethod, error) {
	stripe.Key = c.Secret

	pm, err := paymentmethod.Get(id, nil)
	if err != nil {
		return nil, err
	}
	return pm, nil
}

func (c *Card) GetPaymentIntent(id string) (*stripe.PaymentIntent, error) {
	stripe.Key = c.Secret

	pi, err := paymentintent.Get(id, nil)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		var stripeErr *stripe.Error
		msg := ""
		if errors.As(err, &stripeErr) {
			msg = stripeErr.Msg
		}
		return nil, msg, err
	}

	return pi, "", nil
}

func (c *Card) CreateCustomer(pm, email string) (*stripe.Customer, string, error) {
	stripe.Key = c.Secret
	customerParams := &stripe.CustomerParams{
		PaymentMethod: stripe.String(pm),
		Email:         stripe.String(email),
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(pm),
		},
	}

	cust, err := customer.New(customerParams)
	if err != nil {
		var stripeError *stripe.Error
		msg := ""
		if errors.As(err, &stripeError) {
			msg = stripeError.Msg
		}
		return nil, msg, err
	}

	return cust, "", nil
}

func (c *Card) SubscribeToPlan(cust *stripe.Customer, plan, email, last4, cardType string) (*stripe.Subscription, error) {
	items := []*stripe.SubscriptionItemsParams{
		{Plan: stripe.String(plan)},
	}
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(cust.ID),
		Items:    items,
	}
	params.AddMetadata("last_four", last4)
	params.AddMetadata("card_type", cardType)
	params.AddMetadata("email", email)
	params.AddExpand("latest_invoice.payment_intent")

	subs, err := subscription.New(params)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (c *Card) Refund(pi string, amount int) error {
	stripe.Key = c.Secret
	amountToRefund := int64(amount)

	refundParams := &stripe.RefundParams{
		Amount:        &amountToRefund,
		PaymentIntent: &pi,
	}
	_, err := refund.New(refundParams)
	if err != nil {
		return err
	}

	return nil
}

func (c *Card) CancelSubscription(subscriptionID string) error {
	stripe.Key = c.Secret

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	_, err := subscription.Update(subscriptionID, params)
	if err != nil {
		return err
	}

	return nil
}
