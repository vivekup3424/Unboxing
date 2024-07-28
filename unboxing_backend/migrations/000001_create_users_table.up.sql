CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('Sales', 'Accountant', 'HR', 'Administrator')),
    version INT NOT NULL DEFAULT 1
);
