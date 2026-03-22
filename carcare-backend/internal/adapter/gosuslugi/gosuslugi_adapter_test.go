package gosuslugi

import (
	"context"
	"testing"
)

func TestGosuslugiAdapter_CheckFine(t *testing.T) {
	adapter := NewGosuslugiAdapter()
	ok, err := adapter.CheckFine(context.Background(), "A123BC77")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected no fine, got fine")
	}
}
