package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Customer struct {
	ID        int64     `json:"id"`       // Unique integer ID for each customer
	CreatedAt time.Time `json:"-"`        // Timestamp created for customer automatically when added to the database
	Name      string    `json:"name"`     // Customer's name
	Email     string    `json:"email"`    // Customer's email address
	Phone     string    `json:"phone"`    // Customer's phone number
	Address   string    `json:"address"`  // Customer's address
	Version   int32     `json:"version"`  // Version number for optimistic locking
}

type CustomerModel struct {
	DB *sql.DB
}

// GetAll fetches all customers from the database.
func (m CustomerModel) GetAll() ([]*Customer, error) {
	query := `
	SELECT id, created_at, name, email, phone, address, version
	FROM customers
	ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error getting customers", err)
		return nil, err
	}
	defer rows.Close()

	customers := []*Customer{}

	for rows.Next() {
		var customer Customer

		err = rows.Scan(
			&customer.ID,
			&customer.CreatedAt,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.Address,
			&customer.Version,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

// Insert adds a new customer to the database.
func (m CustomerModel) Insert(customer *Customer) error {
	query := `
	INSERT INTO customers (name, email, phone, address)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	err := m.DB.QueryRow(query, customer.Name, customer.Email, customer.Phone, customer.Address).Scan(&customer.ID, &customer.CreatedAt, &customer.Version)
	if err != nil {
		log.Println("Creating customer in the database", err)
	} else {
		log.Printf("Customer with ID: %d created successfully in the database\n", customer.ID)
	}
	return err
}

// Get fetches a specific customer from the database by ID.
func (m CustomerModel) Get(id int64) (*Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	query := `
	SELECT id, created_at, name, email, phone, address, version 
	FROM customers
	WHERE id = $1
	`

	var customer Customer
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&customer.ID,
		&customer.CreatedAt,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.Version,
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
	return &customer, nil
}

// Update modifies an existing customer in the database.
func (m CustomerModel) Update(customer *Customer) error {
	query := `
	UPDATE customers
	SET name = $1, email = $2, phone = $3, address = $4, version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version
	`

	err := m.DB.QueryRow(query, customer.Name, customer.Email, customer.Phone, customer.Address, customer.ID, customer.Version).Scan(&customer.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			log.Println("Edit conflict (version)", err)
			return ErrEditConflict
		default:
			log.Println("Updating customer", err)
			return err
		}
	} else {
		log.Println("Customer updated successfully")
	}
	return nil
}

// Delete removes a customer from the database.
func (m CustomerModel) Delete(id int64) error {
	query := `
	DELETE FROM customers
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
