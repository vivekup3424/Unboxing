package config

import (
	"flag"
	"os"
)

type DatabaseConfig struct {
	DSN string
}

type Config struct {
	Port     int
	Env      string
	Database DatabaseConfig
}

func LoadConfig() Config {
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.Database.DSN, "db-dsn", os.Getenv("COMPANY_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	return cfg
}
