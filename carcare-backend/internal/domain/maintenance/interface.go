package maintenance

type Repository interface {
	AddMaintenanceEvent(event MaintenanceEvent) error
	GetMaintenanceEvent(id string) (MaintenanceEvent, error)
	UpdateMaintenanceEvent(event MaintenanceEvent) error
	DeleteMaintenanceEvent(id string) error
	ListMaintenanceEvents(userID string) ([]MaintenanceEvent, error)
}
