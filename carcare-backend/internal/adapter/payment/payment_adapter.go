package payment

import (
	"context"
)

// PaymentAdapter defines integration with external payment gateways.
type PaymentAdapter struct {
	// add config, http client, etc.
}

func NewPaymentAdapter(/* config, deps */) *PaymentAdapter {
	return &PaymentAdapter{}
}

// Example: PayFine processes a fine payment.
func (a *PaymentAdapter) PayFine(ctx context.Context, fineID string, amount float64) (string, error) {
	// TODO: Implement real HTTP call to payment gateway
	return "", nil
}
