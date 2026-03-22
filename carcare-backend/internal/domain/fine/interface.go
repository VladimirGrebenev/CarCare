package fine

type Repository interface {
	AddFine(fine Fine) error
}
