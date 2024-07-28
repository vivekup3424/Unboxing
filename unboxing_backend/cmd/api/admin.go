package main

import (
	"company/internal/data"
	"html/template"
	"net/http"
)

func (app *application) showRegisterAdminForm(w http.ResponseWriter, r *http.Request) {
	// Define the template files needed
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/register_admin.tmpl",
	}

	// Parse the template files
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLogger.Println("Error parsing templates:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLogger.Println("Error executing template:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Handle the admin registration form submission
func (app *application) registerAdminHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to process form data", http.StatusBadRequest)
		return
	}

	// Retrieve form values
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	secretKey := r.FormValue("secret-key")
	// Validate form inputs
	if name == "" || email == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if !validateSuperSecretKey(secretKey) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Create a new admin user
	newUser := data.User{
		Name:  name,
		Email: email,
		Role:  "Administrator",
	}
	newUser.Password.Set(password)
	// Insert the new user into the database
	if err := app.models.Users.Insert(&newUser); err != nil {
		http.Error(w, "Database Error in admin", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Admin user created successfully!"))
}

// Helper function to validate the super secret key
func validateSuperSecretKey(authHeader string) bool {
	//if authHeader == "" {
	//	return false
	//}
	//// Extract the key part from the header (assuming "Bearer <key>")
	//parts := strings.Split(authHeader, " ")
	//if len(parts) != 2 || parts[0] != "Bearer" {
	//	return false
	//}
	//providedKey := parts[1]
	//superSecretKey := os.Getenv("SECRET_KEY")
	// Compare the provided key with the expected hash
	return true
}
