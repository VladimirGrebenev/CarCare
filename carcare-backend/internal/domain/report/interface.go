package report

type Repository interface {
	GetByID(id string) (*Report, error)
	// Add more methods as needed
}
