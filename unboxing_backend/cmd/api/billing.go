package main

import (
	"company/internal/data"
	"company/internal/validator"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Date struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Date) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (d Date) MarshalJSON() ([]byte, error) {
	str := d.Format("2006-01-02")
	return json.Marshal(str)
}

func (app *application) showBillingHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Can't get ID (int)", err)
		http.Error(w, "Can't get ID", http.StatusBadRequest)
		return
	}
	billing, err := app.models.Billing.Get(int64(numID))
	if err != nil {
		app.errorLogger.Println("Unable to get billing of this ID", err)
		http.Error(w, "Billing not found", http.StatusNotFound)
		return
	}
	data, err := json.Marshal(billing)
	if err != nil {
		app.errorLogger.Println("Billing data marshalling:", err)
		http.Error(w, "Internal Server Error when getting the billing data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (app *application) createBillingHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CustomerID int64   `json:"customer_id"`
		Amount     float64 `json:"amount"`
		Date       Date    `json:"date"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		app.errorLogger.Println("Decoding request body", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	newBilling := data.Billing{
		CustomerID: input.CustomerID,
		Amount:     input.Amount,
		Date:       input.Date.Time,
	}
	if err := app.models.Billing.Insert(&newBilling); err != nil {
		app.errorLogger.Println("Inserting billing into database", err)
		http.Error(w, "Database Insertion Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(newBilling)
	if err != nil {
		app.errorLogger.Println("Converting the billing struct to JSON", err)
		w.Write([]byte(`"message":"Billing marshalling to JSON failed"`))
	}
	w.Write(js)
	w.Write([]byte(`"message":"New billing created"`))
}

func (app *application) updateBillingHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Can't get ID (int)", err)
		http.Error(w, "Can't get ID", http.StatusBadRequest)
		return
	}
	billing, err := app.models.Billing.Get(int64(numID))
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.errorLogger.Println("Getting billing", err)
			http.Error(w, "Billing not found", http.StatusNotFound)
		} else {
			app.errorLogger.Println("Unknown error getting billing", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	var input struct {
		CustomerID *int64   `json:"customer_id"`
		Amount     *float64 `json:"amount"`
		Date       *Date    `json:"date"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&input)
	if err != nil {
		app.errorLogger.Println("Decoding request body", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if input.CustomerID != nil {
		billing.CustomerID = *input.CustomerID
	}
	if input.Amount != nil {
		billing.Amount = *input.Amount
	}
	if input.Date != nil {
		billing.Date = input.Date.Time
	}
	err = app.models.Billing.Update(billing)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.errorLogger.Println("Edit conflict", err)
			http.Error(w, "Unable to update the record due to edit conflict, please try again", http.StatusConflict)
			return
		default:
			app.errorLogger.Println("Updating billing ID=", billing.ID, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Billing updated successfully"))
	fmt.Fprintf(w, "%+v", billing)
}

func (app *application) deleteBillingHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Can't get ID (int)", err)
		http.Error(w, "Can't get ID", http.StatusBadRequest)
		return
	}
	err = app.models.Billing.Delete(int64(numID))
	if err == data.ErrRecordNotFound {
		app.errorLogger.Println("Billing ID not found", err)
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	} else if err != nil {
		app.errorLogger.Println("Failed delete operation", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	app.infoLogger.Printf("Billing with ID: %v deleted successfully from database\n", numID)
	w.Write([]byte("Billing deleted successfully"))
}

func (app *application) listBillingsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CustomerID int64
		Filters    data.Filters
	}
	queryString := r.URL.Query()

	var customerID int64
	customerIDFromQuery := queryString.Get("customer_id")
	if customerIDFromQuery != "" {
		id, err := strconv.ParseInt(customerIDFromQuery, 10, 64)
		if err == nil {
			customerID = id
		}
	}
	input.CustomerID = customerID
	v := &validator.Validator{}
	input.Filters.Page = app.readInt(queryString, "page", 1, v)
	input.Filters.PageSize = app.readInt(queryString, "page_size", 20, v)
	input.Filters.Sort = app.readString(queryString, "sort", "id")

	billings, err := app.models.Billing.GetAll()
	if err != nil {
		app.errorLogger.Println("Getting billings", err)
		http.Error(w, "Error when getting billings", http.StatusInternalServerError)
		return
	}
	js, err := json.MarshalIndent(billings, "", "\t")
	if err != nil {
		app.errorLogger.Println("Marshalling billing, converting to JSON", err)
		http.Error(w, "Error when getting billings", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(js)
}
