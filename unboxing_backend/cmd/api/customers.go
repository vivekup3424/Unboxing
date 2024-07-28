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

func (app *application) showCustomerHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Can't get ID (int)", err)
		http.Error(w, "Can't get ID", http.StatusBadRequest)
		return
	}
	customer, err := app.models.Customers.Get(int64(numID))
	if err != nil {
		app.errorLogger.Println("Unable to get customer of this ID", err)
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}
	data, err := json.Marshal(customer)
	if err != nil {
		app.errorLogger.Println("Customer data marshalling:", err)
		http.Error(w, "Internal Server Error when getting the customer data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (app *application) createCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"info"`
		Address string `json:"address"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		app.errorLogger.Println("Decoding request body", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	newCustomer := data.Customer{
		Name:    input.Name,
		Email:   input.Email,
		Phone:   input.Phone,
		Address: input.Address,
	}
	if err := app.models.Customers.Insert(&newCustomer); err != nil {
		app.errorLogger.Println("Inserting customer into database", err)
		http.Error(w, "Database Insertion Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(newCustomer)
	if err != nil {
		app.errorLogger.Println("Converting the customer struct to JSON", err)
		w.Write([]byte(`"message":"Customer marshalling to JSON failed"`))
	}
	w.Write(js)
	w.Write([]byte(`"message":"New customer created"`))
}

func (app *application) updateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Can't get ID (int)", err)
		http.Error(w, "Can't get ID", http.StatusBadRequest)
		return
	}
	customer, err := app.models.Customers.Get(int64(numID))
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.errorLogger.Println("Getting customer", err)
			http.Error(w, "Customer not found", http.StatusNotFound)
		} else {
			app.errorLogger.Println("Unknown error getting customer", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	var input struct {
		Name    *string `json:"name"`
		Email   *string `json:"email"`
		Phone   *string `json:"phone"`
		Address *string `json:"address"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&input)
	if err != nil {
		app.errorLogger.Println("Decoding request body", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if input.Name != nil {
		customer.Name = *input.Name
	}
	if input.Email != nil {
		customer.Email = *input.Email
	}
	if input.Phone != nil {
		customer.Phone = *input.Phone
	}
	if input.Address != nil {
		customer.Address = *input.Address
	}
	err = app.models.Customers.Update(customer)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.errorLogger.Println("Edit conflict", err)
			http.Error(w, "Unable to update the record due to edit conflict, please try again", http.StatusConflict)
			return
		default:
			app.errorLogger.Println("Updating customer ID=", customer.ID, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Customer updated successfully"))
	fmt.Fprintf(w, "%+v", customer)
}

func (app *application) deleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Can't get ID (int)", err)
		http.Error(w, "Can't get ID", http.StatusBadRequest)
		return
	}
	err = app.models.Customers.Delete(int64(numID))
	if err == data.ErrRecordNotFound {
		app.errorLogger.Println("Customer ID not found", err)
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	} else if err != nil {
		app.errorLogger.Println("Failed delete operation", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	app.infoLogger.Printf("Customer with ID: %v deleted successfully from database\n", numID)
	w.Write([]byte("Customer deleted successfully"))
}

func (app *application) listCustomersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string
		Filters data.Filters
	}
	queryString := r.URL.Query()

	var name string
	defaultName := ""
	nameFromQuery := queryString.Get("name")
	if nameFromQuery == "" {
		name = defaultName
	} else {
		name = nameFromQuery
	}
	input.Name = name
	v := &validator.Validator{}
	input.Filters.Page = app.readInt(queryString, "page", 1, v)
	input.Filters.PageSize = app.readInt(queryString, "page_size", 20, v)
	input.Filters.Sort = app.readString(queryString, "sort", "id")

	customers, err := app.models.Customers.GetAll()
	if err != nil {
		app.errorLogger.Println("Getting customers", err)
		http.Error(w, "Error when getting customers", http.StatusInternalServerError)
		return
	}
	js, err := json.MarshalIndent(customers, "", "\t")
	if err != nil {
		app.errorLogger.Println("Marshalling customer, converting to JSON", err)
		http.Error(w, "Error when getting customers", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(js)
}
