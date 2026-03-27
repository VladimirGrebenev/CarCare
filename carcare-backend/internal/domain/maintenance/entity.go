package maintenance

type MaintenanceEvent struct {
	ID    string  `json:"id"`
	CarID string  `json:"car_id"`
	Type  string  `json:"type"`
	Date  string  `json:"date"`
	Cost  float64 `json:"cost"`
}
