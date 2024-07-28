CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Insert roles
INSERT INTO roles (name) VALUES 
('Administrator'), 
('HR'), 
('Sales'), 
('Accountant'), 
('Customer');

-- Insert permissions
INSERT INTO permissions (name) VALUES
('manage_employee'), 
('manage_payroll'), 
('view_payroll'), 
('manage_customers'), 
('view_customers'), 
('manage_billing'), 
('view_billing');

-- Assign permissions to roles
INSERT INTO role_permissions (role_id, permission_id) VALUES
    ((SELECT id FROM roles WHERE name = 'Administrator'), (SELECT id FROM permissions WHERE name = 'manage_employee')),
    ((SELECT id FROM roles WHERE name = 'Sales'), (SELECT id FROM permissions WHERE name = 'manage_customers')),
    ((SELECT id FROM roles WHERE name = 'Sales'), (SELECT id FROM permissions WHERE name = 'manage_billing')),
    ((SELECT id FROM roles WHERE name = 'Sales'), (SELECT id FROM permissions WHERE name = 'view_billing')),
    ((SELECT id FROM roles WHERE name = 'HR'), (SELECT id FROM permissions WHERE name = 'manage_payroll')),
    ((SELECT id FROM roles WHERE name = 'HR'), (SELECT id FROM permissions WHERE name = 'view_payroll')),
    ((SELECT id FROM roles WHERE name = 'Accountant'), (SELECT id FROM permissions WHERE name = 'view_payroll')),
    ((SELECT id FROM roles WHERE name = 'Accountant'), (SELECT id FROM permissions WHERE name = 'view_billing'));