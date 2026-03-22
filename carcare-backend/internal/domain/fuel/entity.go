package fuel

type FuelEvent struct {
	ID      string  `json:"id"`
	CarID   string  `json:"car_id"`
	Volume  float64 `json:"volume"`
	Price   float64 `json:"price"`
	Type    string  `json:"type"`
	Date    string  `json:"date"`
}
