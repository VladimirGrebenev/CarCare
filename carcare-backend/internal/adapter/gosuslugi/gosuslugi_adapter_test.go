package gosuslugi

import (
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

// Проверяем, что GosuslugiAdapter реализует интерфейс usecase.GosuslugiFetcher
func TestGosuslugiAdapter_ImplementsGosuslugiFetcher(t *testing.T) {
	var _ usecase.GosuslugiFetcher = NewGosuslugiAdapter()
}
