package fine

type Fine struct {
	ID          string  `json:"id"`
	CarID       string  `json:"car_id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Date        string  `json:"date"`
	Status      string  `json:"status"`
	Description string  `json:"description"`
}
