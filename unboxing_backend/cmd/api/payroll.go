package main

import (
	"company/internal/data"
	"company/internal/validator"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) showPayrollHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Can't get ID (int)", err)
		http.Error(w, "Can't get ID", http.StatusBadRequest)
		return
	}
	payroll, err := app.models.Payroll.Get(int64(numID))
	if err != nil {
		app.errorLogger.Println("Unable to get payroll of this ID", err)
		http.Error(w, "Payroll not found", http.StatusNotFound)
		return
	}
	data, err := json.Marshal(payroll)
	if err != nil {
		app.errorLogger.Println("Payroll data marshalling:", err)
		http.Error(w, "Internal Server Error when getting the payroll data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (app *application) createPayrollHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		EmployeeID int64   `json:"employee_id"`
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
	newPayroll := data.Payroll{
		EmployeeID: input.EmployeeID,
		Amount:     input.Amount,
		Date:       input.Date.Time,
	}
	if err := app.models.Payroll.Insert(&newPayroll); err != nil {
		app.errorLogger.Println("Inserting payroll into database", err)
		http.Error(w, "Database Insertion Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(newPayroll)
	if err != nil {
		app.errorLogger.Println("Converting the payroll struct to JSON", err)
		w.Write([]byte(`"message":"Payroll marshalling to JSON failed"`))
	}
	w.Write(js)
	w.Write([]byte(`"message":"New payroll created"`))
}

func (app *application) updatePayrollHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Can't get ID (int)", err)
		http.Error(w, "Can't get ID", http.StatusBadRequest)
		return
	}
	payroll, err := app.models.Payroll.Get(int64(numID))
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.errorLogger.Println("Getting payroll", err)
			http.Error(w, "Payroll not found", http.StatusNotFound)
		} else {
			app.errorLogger.Println("Unknown error getting payroll", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	var input struct {
		EmployeeID *int64   `json:"employee_id"`
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

	if input.EmployeeID != nil {
		payroll.EmployeeID = *input.EmployeeID
	}
	if input.Amount != nil {
		payroll.Amount = *input.Amount
	}
	if input.Date != nil {
		payroll.Date = *&input.Date.Time
	}
	err = app.models.Payroll.Update(payroll)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.errorLogger.Println("Edit conflict", err)
			http.Error(w, "Unable to update the record due to edit conflict, please try again", http.StatusConflict)
			return
		default:
			app.errorLogger.Println("Updating payroll ID=", payroll.ID, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payroll updated successfully"))
	fmt.Fprintf(w, "%+v", payroll)
}

func (app *application) deletePayrollHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Can't get ID (int)", err)
		http.Error(w, "Can't get ID", http.StatusBadRequest)
		return
	}
	err = app.models.Payroll.Delete(int64(numID))
	if err == data.ErrRecordNotFound {
		app.errorLogger.Println("Payroll ID not found", err)
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	} else if err != nil {
		app.errorLogger.Println("Failed delete operation", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	app.infoLogger.Printf("Payroll with ID: %v deleted successfully from database\n", numID)
	w.Write([]byte("Payroll deleted successfully"))
}

func (app *application) listPayrollsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		EmployeeID int64
		Filters    data.Filters
	}
	queryString := r.URL.Query()

	var employeeID int64
	employeeIDFromQuery := queryString.Get("employee_id")
	if employeeIDFromQuery != "" {
		id, err := strconv.ParseInt(employeeIDFromQuery, 10, 64)
		if err == nil {
			employeeID = id
		}
	}
	input.EmployeeID = employeeID
	v := &validator.Validator{}
	input.Filters.Page = app.readInt(queryString, "page", 1, v)
	input.Filters.PageSize = app.readInt(queryString, "page_size", 20, v)
	input.Filters.Sort = app.readString(queryString, "sort", "id")

	payrolls, err := app.models.Payroll.GetAll()
	if err != nil {
		app.errorLogger.Println("Getting payrolls", err)
		http.Error(w, "Error when getting payrolls", http.StatusInternalServerError)
		return
	}
	js, err := json.MarshalIndent(payrolls, "", "\t")
	if err != nil {
		app.errorLogger.Println("Marshalling payroll, converting to JSON", err)
		http.Error(w, "Error when getting payrolls", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(js)
}
