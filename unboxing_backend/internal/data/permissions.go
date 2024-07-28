package data

import (
	"context"
	"database/sql"
	"time"
)

type Permissions []string

// Add a helper method to check whether the Permissions slice contains a specific permission code.
func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

type PermissionModel struct {
	DB *sql.DB
}

// The GetAllForRole() method returns all permission codes for a role
func (m PermissionModel) GetAllForRole(roleName string) (Permissions, error) {
	query := `
	WITH role_id AS (
	    SELECT id
	    FROM roles
	    WHERE name = $1
	)
	SELECT permissions.name
	FROM permissions
	INNER JOIN role_permissions ON role_permissions.permission_id = permissions.id
	INNER JOIN role_id ON role_permissions.role_id = role_id.id;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var permissions Permissions
	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}
