// FineRepository implements fine.Repository
type FineRepository struct {
	db *sql.DB
}

func NewFineRepository(db *sql.DB) *FineRepository {
	return &FineRepository{db: db}
}

func (r *FineRepository) AddFine(f fine.Fine) error {
	_, err := r.db.Exec(`INSERT INTO fines (id, car_id, amount, type, date, status, description, bill_number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		f.ID, f.CarID, f.Amount, f.Type, f.Date, f.Status, f.Description, f.BillNumber)
	return err
}

func (r *FineRepository) GetFine(id string) (fine.Fine, error) {
	var f fine.Fine
	err := r.db.QueryRow(`SELECT id, car_id, amount, type, date, status, description, COALESCE(bill_number,'') FROM fines WHERE id = $1`, id).
		Scan(&f.ID, &f.CarID, &f.Amount, &f.Type, &f.Date, &f.Status, &f.Description, &f.BillNumber)
	if err == sql.ErrNoRows {
		return fine.Fine{}, errors.New("fine not found")
	}
	return f, err
}

func (r *FineRepository) UpdateFine(f fine.Fine) error {
	res, err := r.db.Exec(`UPDATE fines SET car_id=$1, amount=$2, type=$3, date=$4, status=$5, description=$6, bill_number=$7 WHERE id=$8`,
		f.CarID, f.Amount, f.Type, f.Date, f.Status, f.Description, f.BillNumber, f.ID)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("fine not found")
	}
	return nil
}

func (r *FineRepository) DeleteFine(id string) error {
	res, err := r.db.Exec(`DELETE FROM fines WHERE id=$1`, id)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("fine not found")
	}
	return nil
}

func (r *FineRepository) ListFines(userID string) ([]fine.Fine, error) {
	rows, err := r.db.Query(`SELECT f.id, f.car_id, f.amount, f.type, f.date, f.status, f.description, COALESCE(f.bill_number,'') FROM fines f JOIN cars c ON c.id = f.car_id WHERE c.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fines := make([]fine.Fine, 0)
	for rows.Next() {
		var f fine.Fine
		if err := rows.Scan(&f.ID, &f.CarID, &f.Amount, &f.Type, &f.Date, &f.Status, &f.Description, &f.BillNumber); err != nil {
			return nil, err
		}
		fines = append(fines, f)
	}
	return fines, nil
}

// CheckFineExistsByBillNumber проверяет, есть ли уже штраф с таким номером постановления для данной машины
func (r *FineRepository) CheckFineExistsByBillNumber(carID string, billNumber string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM fines WHERE car_id = $1 AND bill_number = $2 AND bill_number != ''`, carID, billNumber).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ReportRepository implements report.Repository
type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}