package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a sommething that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users       UserModel
	Customers   CustomerModel
	Payroll     PayrollModel
	Billing     BillingModel
	Token       TokenModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:       UserModel{DB: db},
		Customers:   CustomerModel{DB: db},
		Payroll:     PayrollModel{DB: db},
		Billing:     BillingModel{DB: db},
		Token:       TokenModel{DB: db},
		Permissions: PermissionModel{DB: db},
	}
}
