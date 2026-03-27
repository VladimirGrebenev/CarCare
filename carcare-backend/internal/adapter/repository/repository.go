package repository

import (
	"context"
	"errors"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/car"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fine"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fuel"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/maintenance"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

import (
	"database/sql"
	"fmt"
)

type CarRepository struct {
	db *sql.DB
}

func NewCarRepository(db *sql.DB) *CarRepository {
	return &CarRepository{db: db}
}

func (r *CarRepository) AddCar(c car.Car) error {
	_, err := r.db.Exec(`INSERT INTO cars (id, brand, model, year, vin, plate) VALUES ($1, $2, $3, $4, $5, $6)`,
		c.ID, c.Brand, c.Model, c.Year, c.VIN, c.Plate)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("car with VIN already exists: %w", err)
		}
		return err
	}
	return nil
}

func (r *CarRepository) GetCar(id string) (car.Car, error) {
	var c car.Car
	err := r.db.QueryRow(`SELECT id, brand, model, year, vin, plate FROM cars WHERE id = $1`, id).
		Scan(&c.ID, &c.Brand, &c.Model, &c.Year, &c.VIN, &c.Plate)
	if err == sql.ErrNoRows {
		return car.Car{}, fmt.Errorf("car not found")
	}
	if err != nil {
		return car.Car{}, err
	}
	return c, nil
}

func (r *CarRepository) UpdateCar(c car.Car) error {
	res, err := r.db.Exec(`UPDATE cars SET brand=$1, model=$2, year=$3, vin=$4, plate=$5 WHERE id=$6`,
		c.Brand, c.Model, c.Year, c.VIN, c.Plate, c.ID)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("car not found")
	}
	return nil
}

func (r *CarRepository) DeleteCar(id string) error {
	res, err := r.db.Exec(`DELETE FROM cars WHERE id=$1`, id)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("car not found")
	}
	return nil
}

func (r *CarRepository) ListCars() ([]car.Car, error) {
	rows, err := r.db.Query(`SELECT id, brand, model, year, vin, plate FROM cars`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cars := make([]car.Car, 0)
	for rows.Next() {
		var c car.Car
		if err := rows.Scan(&c.ID, &c.Brand, &c.Model, &c.Year, &c.VIN, &c.Plate); err != nil {
			return nil, err
		}
		cars = append(cars, c)
	}
	return cars, nil
}

// isUniqueViolation checks if error is pq unique violation
func isUniqueViolation(err error) bool {
	return err != nil && (err.Error() == "pq: duplicate key value violates unique constraint \"cars_vin_key\"")
}

// UserRepository implements user.Repository
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	_, err = r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, name, role, password_hash) VALUES ($1, $2, $3, $4, $5)`,
		u.ID, string(u.Email), u.Name, string(u.Role), string(hashed))
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("user with email already exists: %w", err)
		}
		return err
	}
	// Сохраним хеш обратно в объект, чтобы caller получил актуальное значение
	u.PasswordHash = string(hashed)
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, email, name, role, COALESCE(password_hash,'') FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, email, name, role, COALESCE(password_hash,'') FROM users WHERE email = $1`, string(email))
	return scanUser(row)
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	res, err := r.db.ExecContext(ctx,
		`UPDATE users SET email=$1, name=$2, role=$3 WHERE id=$4`,
		string(u.Email), u.Name, string(u.Role), u.ID)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context) ([]*user.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, email, name, role, COALESCE(password_hash,'') FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*user.User, 0)
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func scanUser(scanner interface {
	Scan(dest ...any) error
}) (*user.User, error) {
	var id, email, name, role, passwordHash string
	if err := scanner.Scan(&id, &email, &name, &role, &passwordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	u := &user.User{
		ID:           id,
		Email:        user.Email(email),
		Name:         name,
		Role:         user.Role(role),
		PasswordHash: passwordHash,
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	return u, nil
}

// FuelRepository implements fuel.Repository
type FuelRepository struct {
       db *sql.DB
}

func NewFuelRepository(db *sql.DB) *FuelRepository {
       return &FuelRepository{db: db}
}

func (r *FuelRepository) AddFuelEvent(event fuel.FuelEvent) error {
       _, err := r.db.Exec(`INSERT INTO fuel_events (id, car_id, volume, price, type, date) VALUES ($1, $2, $3, $4, $5, $6)`,
	       event.ID, event.CarID, event.Volume, event.Price, event.Type, event.Date)
       return err
}

func (r *FuelRepository) GetFuelEvent(id string) (fuel.FuelEvent, error) {
       var e fuel.FuelEvent
       err := r.db.QueryRow(`SELECT id, car_id, volume, price, type, date FROM fuel_events WHERE id = $1`, id).
	       Scan(&e.ID, &e.CarID, &e.Volume, &e.Price, &e.Type, &e.Date)
       if err == sql.ErrNoRows {
	       return fuel.FuelEvent{}, errors.New("fuel event not found")
       }
       return e, err
}

func (r *FuelRepository) UpdateFuelEvent(event fuel.FuelEvent) error {
       res, err := r.db.Exec(`UPDATE fuel_events SET car_id=$1, volume=$2, price=$3, type=$4, date=$5 WHERE id=$6`,
	       event.CarID, event.Volume, event.Price, event.Type, event.Date, event.ID)
       if err != nil {
	       return err
       }
       count, _ := res.RowsAffected()
       if count == 0 {
	       return errors.New("fuel event not found")
       }
       return nil
}

func (r *FuelRepository) DeleteFuelEvent(id string) error {
       res, err := r.db.Exec(`DELETE FROM fuel_events WHERE id=$1`, id)
       if err != nil {
	       return err
       }
       count, _ := res.RowsAffected()
       if count == 0 {
	       return errors.New("fuel event not found")
       }
       return nil
}

func (r *FuelRepository) ListFuelEvents() ([]fuel.FuelEvent, error) {
       rows, err := r.db.Query(`SELECT id, car_id, volume, price, type, date FROM fuel_events`)
       if err != nil {
	       return nil, err
       }
       defer rows.Close()
       events := make([]fuel.FuelEvent, 0)
       for rows.Next() {
	       var e fuel.FuelEvent
	       if err := rows.Scan(&e.ID, &e.CarID, &e.Volume, &e.Price, &e.Type, &e.Date); err != nil {
		       return nil, err
	       }
	       events = append(events, e)
       }
       return events, nil
}

// MaintenanceRepository implements maintenance.Repository
type MaintenanceRepository struct {
	db *sql.DB
}

func NewMaintenanceRepository(db *sql.DB) *MaintenanceRepository {
	return &MaintenanceRepository{db: db}
}

func (r *MaintenanceRepository) AddMaintenanceEvent(event maintenance.MaintenanceEvent) error {
	_, err := r.db.Exec(`INSERT INTO maintenance_events (id, car_id, type, date, cost) VALUES ($1, $2, $3, $4, $5)`,
		event.ID, event.CarID, event.Type, event.Date, event.Cost)
	return err
}

func (r *MaintenanceRepository) GetMaintenanceEvent(id string) (maintenance.MaintenanceEvent, error) {
	var e maintenance.MaintenanceEvent
	err := r.db.QueryRow(`SELECT id, car_id, type, date, cost FROM maintenance_events WHERE id = $1`, id).
		Scan(&e.ID, &e.CarID, &e.Type, &e.Date, &e.Cost)
	if err == sql.ErrNoRows {
		return maintenance.MaintenanceEvent{}, errors.New("maintenance event not found")
	}
	return e, err
}

func (r *MaintenanceRepository) UpdateMaintenanceEvent(event maintenance.MaintenanceEvent) error {
	res, err := r.db.Exec(`UPDATE maintenance_events SET car_id=$1, type=$2, date=$3, cost=$4 WHERE id=$5`,
		event.CarID, event.Type, event.Date, event.Cost, event.ID)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("maintenance event not found")
	}
	return nil
}

func (r *MaintenanceRepository) DeleteMaintenanceEvent(id string) error {
	res, err := r.db.Exec(`DELETE FROM maintenance_events WHERE id=$1`, id)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("maintenance event not found")
	}
	return nil
}

func (r *MaintenanceRepository) ListMaintenanceEvents() ([]maintenance.MaintenanceEvent, error) {
	rows, err := r.db.Query(`SELECT id, car_id, type, date, cost FROM maintenance_events`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := make([]maintenance.MaintenanceEvent, 0)
	for rows.Next() {
		var e maintenance.MaintenanceEvent
		if err := rows.Scan(&e.ID, &e.CarID, &e.Type, &e.Date, &e.Cost); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

// FineRepository implements fine.Repository
type FineRepository struct {
	db *sql.DB
}

func NewFineRepository(db *sql.DB) *FineRepository {
	return &FineRepository{db: db}
}

func (r *FineRepository) AddFine(f fine.Fine) error {
	_, err := r.db.Exec(`INSERT INTO fines (id, car_id, amount, type, date, status, description) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		f.ID, f.CarID, f.Amount, f.Type, f.Date, f.Status, f.Description)
	return err
}

func (r *FineRepository) GetFine(id string) (fine.Fine, error) {
	var f fine.Fine
	err := r.db.QueryRow(`SELECT id, car_id, amount, type, date, status, description FROM fines WHERE id = $1`, id).
		Scan(&f.ID, &f.CarID, &f.Amount, &f.Type, &f.Date, &f.Status, &f.Description)
	if err == sql.ErrNoRows {
		return fine.Fine{}, errors.New("fine not found")
	}
	return f, err
}

func (r *FineRepository) UpdateFine(f fine.Fine) error {
	res, err := r.db.Exec(`UPDATE fines SET car_id=$1, amount=$2, type=$3, date=$4, status=$5, description=$6 WHERE id=$7`,
		f.CarID, f.Amount, f.Type, f.Date, f.Status, f.Description, f.ID)
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

func (r *FineRepository) ListFines() ([]fine.Fine, error) {
	rows, err := r.db.Query(`SELECT id, car_id, amount, type, date, status, description FROM fines`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fines := make([]fine.Fine, 0)
	for rows.Next() {
		var f fine.Fine
		if err := rows.Scan(&f.ID, &f.CarID, &f.Amount, &f.Type, &f.Date, &f.Status, &f.Description); err != nil {
			return nil, err
		}
		fines = append(fines, f)
	}
	return fines, nil
}

// ReportRepository implements report.Repository
type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}
