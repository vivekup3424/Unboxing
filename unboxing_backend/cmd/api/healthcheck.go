package main

import (
	"fmt"
	"net/http"
)

// HealthCheckData holds data for the health check page
type HealthCheckData struct {
	Status      string
	Environment string
	Version     string
}

// healthcheckHandler writes a plain-text response with information about the application status, operating environment, and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Respond with status information (this part is now redundant as the template renders the info)
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}
