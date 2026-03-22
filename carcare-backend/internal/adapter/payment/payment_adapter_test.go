package payment

import (
	"context"
	"testing"
)

func TestPaymentAdapter_PayFine(t *testing.T) {
	adapter := NewPaymentAdapter()
	id, err := adapter.PayFine(context.Background(), "fine123", 100.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "" {
		t.Errorf("expected empty id, got %v", id)
	}
}
