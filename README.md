# Company Management System

This is a demo application for managing a company's operations using a web-based interface. The application is built using Go for the backend and HTML, CSS, and JavaScript for the frontend. It supports functionalities for user management, customer management, billing, payroll, and role-based access control with different permissions.

https://github.com/user-attachments/assets/740361c4-3639-4d9c-9a38-18106eae4a0b

## Features

- **Role-Based Access Control (RBAC)**: Implements roles such as Administrator, HR, Sales, and Accountant, each with specific permissions. This ensures that users have access only to the parts of the application they are authorized to manage.
- **User Management**: Administrators can manage users, including creating, updating, and deleting user accounts.

- **Customer Management**: Sales personnel can manage customer information.

- **Billing Management**: Manage billing information with permissions for viewing and editing.

- **Payroll Management**: HR and Accountants can view payroll data, while HR can edit it.

- **Stateful Token-Based Authentication**: Uses JWT tokens to authenticate users and manage sessions securely.

## Architecture

![diagram](https://github.com/user-attachments/assets/9b2ced64-16b6-4b7e-a79b-638e3ad7ac2f)

- **Backend**: Go (Golang) for server-side application logic.
- **Frontend**: HTML, CSS, and JavaScript for the user interface.
- **Database**: PostgreSQL for data storage.

## Middlewares Used

- **Security Headers**: Adds security-related headers to HTTP responses.
- **CORS**: Enables Cross-Origin Resource Sharing.
- **Authentication**: Handles token-based authentication.
- **Permission Checking**: Ensures users have the required permissions to access specific endpoints.

## Context Usage in Go

- **Context**: Used for managing request-scoped values, deadlines, and cancellations.

## Technologies Used

- **Go (Golang)**
- **HTML**
- **CSS**
- **JavaScript**
- **PostgreSQL**
  
## Lines of Code
- Backend Golang Application with Migrations
  ![lines](https://github.com/user-attachments/assets/d5ea60a6-1edf-42ed-b3cd-d82aa4cd9880)
- 
---

This description should give a clear and concise overview of your application.
