package fine

type Repository interface {
	AddFine(fine Fine) error
	GetFine(id string) (Fine, error)
	UpdateFine(fine Fine) error
	DeleteFine(id string) error
	ListFines() ([]Fine, error)
}
