package maintenance

type Repository interface {
	AddMaintenanceEvent(event MaintenanceEvent) error
}
