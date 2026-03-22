package gosuslugi

import (
	"context"
)

// GosuslugiAdapter defines integration with the Gosuslugi external service.
type GosuslugiAdapter struct {
	// add config, http client, etc.
}

func NewGosuslugiAdapter(/* config, deps */) *GosuslugiAdapter {
	return &GosuslugiAdapter{}
}

// Example: CheckFine queries fines by car number or user passport.
func (a *GosuslugiAdapter) CheckFine(ctx context.Context, carNumber string) (bool, error) {
	// TODO: Implement real HTTP call to Gosuslugi API
	return false, nil
}
