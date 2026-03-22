package maintenance

type MaintenanceEvent struct {
	ID     string
	CarID  string
	Type   string
	Date   string
	Cost   float64
}
