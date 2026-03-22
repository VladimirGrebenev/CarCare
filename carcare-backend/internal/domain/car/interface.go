package car

type Repository interface {
	AddCar(car Car) error
	GetCar(id string) (Car, error)
	UpdateCar(car Car) error
	DeleteCar(id string) error
	ListCars() ([]Car, error)
}
