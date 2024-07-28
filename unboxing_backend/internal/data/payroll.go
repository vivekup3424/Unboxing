package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type Payroll struct {
	ID        int64     `json:"id"`        // Unique integer ID for each payroll entry
	EmployeeID int64    `json:"employee_id"`// Employee ID to whom the payroll belongs
	Amount    float64   `json:"amount"`    // Payroll amount
	Date      time.Time `json:"date"`      // Payroll date
	Version   int32     `json:"version"`   // Version number for optimistic locking
}

type PayrollModel struct {
	DB *sql.DB
}

// GetAll fetches all payroll entries from the database.
func (m PayrollModel) GetAll() ([]*Payroll, error) {
	query := `
	SELECT id, employee_id, amount, date, version
	FROM payroll
	ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error getting payroll entries", err)
		return nil, err
	}
	defer rows.Close()

	payrolls := []*Payroll{}

	for rows.Next() {
		var payroll Payroll

		err = rows.Scan(
			&payroll.ID,
			&payroll.EmployeeID,
			&payroll.Amount,
			&payroll.Date,
			&payroll.Version,
		)
		if err != nil {
			return nil, err
		}
		payrolls = append(payrolls, &payroll)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payrolls, nil
}

// Insert adds a new payroll entry to the database.
func (m PayrollModel) Insert(payroll *Payroll) error {
	query := `
	INSERT INTO payroll (employee_id, amount, date)
	VALUES ($1, $2, $3)
	RETURNING id, version
	`

	err := m.DB.QueryRow(query, payroll.EmployeeID, payroll.Amount, payroll.Date).Scan(&payroll.ID, &payroll.Version)
	if err != nil {
		log.Println("Creating payroll entry in the database", err)
	} else {
		log.Printf("Payroll entry with ID: %d created successfully in the database\n", payroll.ID)
	}
	return err
}

// Get fetches a specific payroll entry from the database by ID.
func (m PayrollModel) Get(id int64) (*Payroll, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	query := `
	SELECT id, employee_id, amount, date, version 
	FROM payroll
	WHERE id = $1
	`

	var payroll Payroll
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&payroll.ID,
		&payroll.EmployeeID,
		&payroll.Amount,
		&payroll.Date,
		&payroll.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Data row not found in the database", err)
			return nil, ErrRecordNotFound
		} else {
			log.Println("Unknown error occurred", err)
			return nil, err
		}
	}
	return &payroll, nil
}

// Update modifies an existing payroll entry in the database.
func (m PayrollModel) Update(payroll *Payroll) error {
	query := `
	UPDATE payroll
	SET employee_id = $1, amount = $2, date = $3, version = version + 1
	WHERE id = $4 AND version = $5
	RETURNING version
	`

	err := m.DB.QueryRow(query, payroll.EmployeeID, payroll.Amount, payroll.Date, payroll.ID, payroll.Version).Scan(&payroll.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			log.Println("Edit conflict (version)", err)
			return ErrEditConflict
		default:
			log.Println("Updating payroll entry", err)
			return err
		}
	} else {
		log.Println("Payroll entry updated successfully")
	}
	return nil
}

// Delete removes a payroll entry from the database.
func (m PayrollModel) Delete(id int64) error {
	query := `
	DELETE FROM payroll
	WHERE id = $1
	`

	results, err := m.DB.Exec(query, id)
	if err != nil {
		log.Println("Delete operation", err)
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected", err)
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
