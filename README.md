# Company Management System

This is a demo application for managing a company's operations using a web-based interface. The application is built using Go for the backend and HTML, CSS, and JavaScript for the frontend. It supports functionalities for user management, customer management, billing, payroll, and role-based access control with different permissions.

[![IMAGE ALT TEXT HERE](https://img.youtube.com/vi/YOUTUBE_VIDEO_ID_HERE/0.jpg)](https://www.youtube.com/watch?v=Vrau6EMr8eo)
## Features

- **Role-Based Access Control (RBAC)**: Implements roles such as Administrator, HR, Sales, and Accountant, each with specific permissions. This ensures that users have access only to the parts of the application they are authorized to manage.
- **User Management**: Administrators can manage users, including creating, updating, and deleting user accounts.

- **Customer Management**: Sales personnel can manage customer information.

- **Billing Management**: Manage billing information with permissions for viewing and editing.

- **Payroll Management**: HR and Accountants can view payroll data, while HR can edit it.

- **Stateful Token-Based Authentication**: Uses JWT tokens to authenticate users and manage sessions securely.

## Architecture

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

---

This description should give a clear and concise overview of your application.
