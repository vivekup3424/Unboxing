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
)

// Enum Role
type Role int

const (
	Administrator Role = iota
	Sales
	Accountant
	HR
)

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Can't get ID (int)", err)
		http.Error(w, "Can't get ID", http.StatusBadRequest)
		return
	}
	user, err := app.models.Users.Get(int64(numID))
	if err != nil {
		app.errorLogger.Println("Unable to get user of this ID", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	data, err := json.Marshal(user)
	if err != nil {
		app.errorLogger.Println("User data marshalling:", err)
		http.Error(w, "Internal Server Error when getting the user data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we
	// expect to be in the HTTP request body
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		app.errorLogger.Println("Decoding request body", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// Feeding the data to the database
	newUser := data.User{
		Name: input.Name,
	}
	if err := app.models.Users.Insert(&newUser); err != nil {
		app.errorLogger.Println("Inserting user into database", err)
		http.Error(w, "Database Insertion Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(newUser)
	if err != nil {
		app.errorLogger.Println("Converting the user struct to JSON", err)
		w.Write([]byte(`"message":"User marshalling to JSON failed"`))
	}
	w.Write(js)
	w.Write([]byte(`"message":"New user created"`))
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Can't get ID (int)", err)
		http.Error(w, "Can't get ID", http.StatusBadRequest)
		return
	}
	user, err := app.models.Users.Get(int64(numID))
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.errorLogger.Println("Getting user", err)
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			app.errorLogger.Println("Unknown error getting user", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Get the new values for user update
	var input struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Role     *string `json:"role"`
		Password *string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&input)
	if err != nil {
		app.errorLogger.Println("Decoding request body", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Copy the values from input to the user pointer
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.Password != nil {
		user.Password.Set(*input.Password)
	}
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.errorLogger.Println("Edit conflict", err)
			http.Error(w, "Unable to update the record due to edit conflict, please try again", http.StatusConflict)
			return
		default:
			app.errorLogger.Println("Updating user ID=", user.ID, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
	fmt.Fprintf(w, "%+v", user)
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Can't get ID (int)", err)
		http.Error(w, "Can't get ID", http.StatusBadRequest)
		return
	}
	err = app.models.Users.Delete(int64(numID))
	if err == data.ErrRecordNotFound {
		app.errorLogger.Println("User ID not found", err)
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	} else if err != nil {
		app.errorLogger.Println("Failed delete operation", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	app.infoLogger.Printf("User with ID: %v deleted successfully from database\n", numID)
	w.Write([]byte("User deleted successfully"))
}

func (app *application) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string
		Roles   []string
		Filters data.Filters
	}
	// Parse and get the params from the query string
	queryString := r.URL.Query()

	var name string
	defaultName := ""
	// Get the name from query string
	nameFromQuery := queryString.Get("name")
	if nameFromQuery == "" {
		name = defaultName
	} else {
		name = nameFromQuery
	}
	input.Name = name

	// Get the roles from the query string
	var roles []string
	defaultRoles := []string{}
	rolesFromQuery := queryString.Get("roles")
	if rolesFromQuery == "" {
		roles = defaultRoles
	} else {
		roles = strings.Split(rolesFromQuery, ",")
	}
	input.Roles = roles
	v := &validator.Validator{}
	// Get the page number
	input.Filters.Page = app.readInt(queryString, "page", 1, v)
	input.Filters.PageSize = app.readInt(queryString, "page_size", 20, v)

	// Extract the sort query string value, falling back to "id" if it is not provided
	input.Filters.Sort = app.readString(queryString, "sort", "id")

	users, err := app.models.Users.GetAll()
	if err != nil {
		app.errorLogger.Println("Getting users", err)
		http.Error(w, "Error when getting users", http.StatusInternalServerError)
		return
	}
	js, err := json.MarshalIndent(users, "", "\t")
	if err != nil {
		app.errorLogger.Println("Marshalling user, converting to JSON", err)
		http.Error(w, "Error when getting users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(js)
}

//GetForToken() returns the user, associated with a token
