package main

import (
	"company/internal/data"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

// Define a config struct to hold all the configuration settings for our application.
// For now, the only configuration settings will be the network port that we want the
// server to listen on, and the name of the current operating environment for the
// application (development, staging, production, etc.). We will read in these
// configuration settings from command-line flags when the application starts.
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config      config
	errorLogger *log.Logger
	infoLogger  *log.Logger
	models      data.Models
}

func main() {
	//declate an instance of config struct
	var cfg config

	//read the value of port and env-commandline flags into the
	//config struct.
	//defaults are 4000 port and "development" environment
	flag.IntVar(&cfg.port, "port", 4000, "API Server Port")
	flag.StringVar(&cfg.env, "env", "development",
		"Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("COMPANY_DB_DSN"), "POSTGRESQL DSN")
	flag.Parse()

	//logger to write message to stdout
	infoLogger := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)
	//an instance of the application struct
	app := application{
		config:      cfg,
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
	}
	// Load templates

	//connect to database, and open a connection
	db, err := app.openDB()
	if err != nil {
		errorLogger.Fatal(err)
	}
	app.models = data.NewModels(db) //is it ok to have a circular dependency here
	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()
	// Also log a message to say that the connection pool has been successfully
	// established.
	infoLogger.Printf("database connection pool established")
	router := http.NewServeMux()
	app.enableCORS(router)
	router.HandleFunc("GET /healthcheck", app.healthcheckHandler)
	router.HandleFunc(("/tokens/authentication"), app.createAuthenticationTokenHandler)
	router.HandleFunc("POST /admin/register", app.registerAdminHandler)

	//user management done by Adminstrator
	router.HandleFunc("GET /v1/user",
		app.requirePermission("manage_employee", app.listUsersHandler))
	router.HandleFunc("POST /v1/user",
		app.requirePermission("manage_employee", app.createUserHandler)) //register new user
	router.HandleFunc("GET /v1/user/{id}",
		app.requirePermission("manage_employee", app.showUserHandler))
	router.HandleFunc("PATCH /v1/user/{id}",
		app.requirePermission("manage_employee", app.updateUserHandler))
	router.HandleFunc("DELETE /v1/user/{id}",
		app.requirePermission("manage_employee", app.deleteUserHandler))
	//customers management done by Sales guy
	router.HandleFunc("GET /v1/customer", app.requirePermission("manage_customers", app.listCustomersHandler))
	router.HandleFunc("POST /v1/customer", app.requirePermission("manage_customers", app.createCustomerHandler))
	router.HandleFunc("GET /v1/customer/{id}", app.requirePermission("manage_customers", app.showCustomerHandler))
	router.HandleFunc("PATCH /v1/customer/{id}", app.requirePermission("manage_customers", app.updateCustomerHandler))
	router.HandleFunc("DELETE /v1/customer/{id}", app.requirePermission("manage_customers", app.deleteCustomerHandler))

	//billing, accountants and sales guy can view it, but only sales guy can change it
	router.HandleFunc("GET /v1/billing", app.requirePermission("view_billing", app.listBillingsHandler))
	router.HandleFunc("POST /v1/billing", app.requirePermission("manage_billing", app.createBillingHandler))
	router.HandleFunc("GET /v1/billing/{id}", app.requirePermission("view_billing", app.showBillingHandler))
	router.HandleFunc("PATCH /v1/billing/{id}", app.requirePermission("manage_billing", app.updateBillingHandler))
	router.HandleFunc("DELETE /v1/billing/{id}", app.requirePermission("manage_billing", app.deleteBillingHandler))

	//payroll similarly accountants and HR can view it, but only HR can change it
	router.HandleFunc("GET /v1/payroll", app.requirePermission("view_payroll", app.listPayrollsHandler))
	router.HandleFunc("POST /v1/payroll", app.requirePermission("manage_payroll", app.createPayrollHandler))
	router.HandleFunc("GET /v1/payroll/{id}", app.requirePermission("view_payroll", app.showPayrollHandler))
	router.HandleFunc("PATCH /v1/payroll/{id}", app.requirePermission("manage_payroll", app.updatePayrollHandler))
	router.HandleFunc("DELETE /v1/payroll/{id}", app.requirePermission("manage_payroll", app.deleteUserHandler))

	//attaching middlewares
	//router.Handle("/v1", app.authenticate(router))
	//declare a http with some good timeout settings. >>>>ich listens
	//on the provided with port, and the above router as the handler
	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", cfg.port),
		Handler:      app.enableCORS(router),
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	//Start the http server
	infoLogger.Printf("starting the %s server on : %s", cfg.env, srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		errorLogger.Fatal(err)
	}
}

// The openDB() function returns a sql.DB connection pool.
// creating an dependency to app
func (app *application) openDB() (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", app.config.db.dsn)
	if err != nil {
		app.errorLogger.Fatal("Connecting to database", err)
		return nil, err
	}
	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		app.errorLogger.Fatal("Pinging to database", err)
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}
