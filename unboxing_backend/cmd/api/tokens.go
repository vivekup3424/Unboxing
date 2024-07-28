package main

import (
	"company/internal/data"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type Envelope map[string]interface{}

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	input.Email = r.FormValue("email")
	input.Password = r.FormValue("password")
	// Debugging: Log the input received
	log.Printf("Received input: %+v\n", input)

	// Lookup the user record based on the email address
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorLogger.Println("Invalid email address", err)
			http.Error(w, "Invalid email address", http.StatusUnauthorized)
		} else {
			app.errorLogger.Println("Querying database", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Validate password
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.errorLogger.Println("Comparing password", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !match {
		app.errorLogger.Println("Invalid password")
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate a new token
	token, err := app.models.Token.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.errorLogger.Println("Creating new token", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Debugging: Log the created token
	log.Printf("Created token: %v\n", token)

	// Encode the token to JSON and send it in the response
	response := Envelope{"authentication_token": token}
	responseJson, err := json.Marshal(response)
	if err != nil {
		app.errorLogger.Println("Error marshaling response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJson)
}
