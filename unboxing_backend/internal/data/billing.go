package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Billing struct {
	ID         int64     `json:"id"`         // Unique integer ID for each billing entry
	CustomerID int64     `json:"customer_id"`// Customer ID to whom the billing belongs
	Amount     float64   `json:"amount"`     // Billing amount
	Date       time.Time `json:"date"`       // Billing date
	Version    int32     `json:"version"`    // Version number for optimistic locking
}

type BillingModel struct {
	DB *sql.DB
}

// GetAll fetches all billing entries from the database.
func (m BillingModel) GetAll() ([]*Billing, error) {
	query := `
	SELECT id, customer_id, amount, date, version
	FROM billing
	ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error getting billing entries", err)
		return nil, err
	}
	defer rows.Close()

	billings := []*Billing{}

	for rows.Next() {
		var billing Billing

		err = rows.Scan(
			&billing.ID,
			&billing.CustomerID,
			&billing.Amount,
			&billing.Date,
			&billing.Version,
		)
		if err != nil {
			return nil, err
		}
		billings = append(billings, &billing)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return billings, nil
}

// Insert adds a new billing entry to the database.
func (m BillingModel) Insert(billing *Billing) error {
	query := `
	INSERT INTO billing (customer_id, amount, date)
	VALUES ($1, $2, $3)
	RETURNING id, version
	`

	err := m.DB.QueryRow(query, billing.CustomerID, billing.Amount, billing.Date).Scan(&billing.ID, &billing.Version)
	if err != nil {
		log.Println("Creating billing entry in the database", err)
	} else {
		log.Printf("Billing entry with ID: %d created successfully in the database\n", billing.ID)
	}
	return err
}

// Get fetches a specific billing entry from the database by ID.
func (m BillingModel) Get(id int64) (*Billing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	query := `
	SELECT id, customer_id, amount, date, version 
	FROM billing
	WHERE id = $1
	`

	var billing Billing
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&billing.ID,
		&billing.CustomerID,
		&billing.Amount,
		&billing.Date,
		&billing.Version,
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
	return &billing, nil
}

// Update modifies an existing billing entry in the database.
func (m BillingModel) Update(billing *Billing) error {
	query := `
	UPDATE billing
	SET customer_id = $1, amount = $2, date = $3, version = version + 1
	WHERE id = $4 AND version = $5
	RETURNING version
	`

	err := m.DB.QueryRow(query, billing.CustomerID, billing.Amount, billing.Date, billing.ID, billing.Version).Scan(&billing.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			log.Println("Edit conflict (version)", err)
			return ErrEditConflict
		default:
			log.Println("Updating billing entry", err)
			return err
		}
	} else {
		log.Println("Billing entry updated successfully")
	}
	return nil
}

// Delete removes a billing entry from the database.
func (m BillingModel) Delete(id int64) error {
	query := `
	DELETE FROM billing
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
