package car

type Repository interface {
	AddCar(car Car) error
	GetCar(id string) (Car, error)
	UpdateCar(car Car, userID string) error
	DeleteCar(id string, userID string) error
	ListCars(userID string) ([]Car, error)
}
