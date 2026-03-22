package maps

import (
	"context"
	"testing"
)

func TestMapsAdapter_GetRoute(t *testing.T) {
	adapter := NewMapsAdapter()
	route, err := adapter.GetRoute(context.Background(), "A", "B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if route != "" {
		t.Errorf("expected empty route, got %v", route)
	}
}
