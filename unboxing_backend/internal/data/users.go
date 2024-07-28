package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// custom error messages
var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type password struct {
	plaintext *string
	hash      []byte
}
type User struct {
	ID        int64     `json:"id"`         // Unique integer ID for each user
	CreatedAt time.Time `json:"created_at"` // Timestamp created for user automatically when added to the database
	Name      string    `json:"name"`       // User's name
	Email     string    `json:"email"`      // User's email address
	Role      string    `json:"role"`       // User's role (Administrator, HR, Sales, Accountant)
	Password  password  `json:"-"`
	Version   int32     `json:"-"` // Version number for optimistic locking
}

// The Set() method calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

// The Matches() method checks whether the provided plaintext password matches the
// hashed password stored in the struct, returning true if it matches and false
// otherwise.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type UserModel struct {
	DB *sql.DB
}

// GetAll fetches all users from the database.
func (m UserModel) GetAll() ([]*User, error) {
	query := `
	SELECT id, created_at, name, email, role, version
	FROM users
	ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error getting users", err)
		return nil, err
	}
	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var user User

		err = rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.Version,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Insert adds a new user to the database.
func (m UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (name, email, password_hash, role)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`
	args := []interface{}{user.Name, user.Email, user.Password.hash, user.Role}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	//using spread operator here
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		log.Println("Error creating user", err)
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	} else {
		log.Printf("User with ID: %d created successfully in the database\n", user.ID)
	}
	return err
}
func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id,created_at,name,email,password_hash,role,version
	FROM users
	WHERE email = $1
	`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Role,
		&user.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}
	return &user, nil
}

// Get fetches a specific user from the database by ID.
func (m UserModel) Get(id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	query := `
	SELECT id, created_at, name, email, role, version 
	FROM users
	WHERE id = $1
	`

	var user User
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.Version,
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
	return &user, nil
}

// Update modifies an existing user in the database.
func (m UserModel) Update(user *User) error {
	query := `
	UPDATE users
	SET name = $1, email = $2,password_hash = $3, role = $4, version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version
	`
	args := []interface{}{user.Name, user.Email, user.Password.hash, user.Role, user.ID, user.Version}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unqiue constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			log.Println("Edit conflict (version)", err)
			return ErrEditConflict
		default:
			log.Println("Error Updating user", err)
			return err
		}
	} else {
		log.Println("User updated successfully")
	}
	return nil
}

// Delete removes a user from the database.
func (m UserModel) Delete(id int64) error {
	query := `
	DELETE FROM users
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

// retrievs the information assosiated with the token for
// a a particular header
func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	query := `
SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.role, users.version
FROM users
INNER JOIN tokens
ON users.id = tokens.user_id
WHERE tokens.hash = $1
AND tokens.scope = $2
AND tokens.expiry > $3`
	args := []interface{}{tokenPlaintext, tokenScope, time.Now()}
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Name, &user.Email, &user.Password.hash, &user.Role, &user.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
