package fuel

type Repository interface {
	AddFuelEvent(event FuelEvent) error
	GetFuelEvent(id string) (FuelEvent, error)
	UpdateFuelEvent(event FuelEvent) error
	DeleteFuelEvent(id string) error
	ListFuelEvents(userID string) ([]FuelEvent, error)
}
