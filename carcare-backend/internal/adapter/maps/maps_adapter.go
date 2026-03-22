package maps

import (
	"context"
)

// MapsAdapter defines integration with external map services (e.g., Yandex, Google).
type MapsAdapter struct {
	// add config, http client, etc.
}

func NewMapsAdapter(/* config, deps */) *MapsAdapter {
	return &MapsAdapter{}
}

// Example: GetRoute calculates a route between two points.
func (a *MapsAdapter) GetRoute(ctx context.Context, from, to string) (string, error) {
	// TODO: Implement real HTTP call to map service
	return "", nil
}
