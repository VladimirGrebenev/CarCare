package car

type Car struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Brand  string `json:"brand"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	VIN    string `json:"vin"`
	Plate  string `json:"plate"`
}
