# CodeBin

A production-ready RESTful API service for creating and sharing code snippets, built with Go and PostgreSQL. Live Demo UI [https://codebin-restful-api-service.onrender.com/](https://codebin-restful-api-service.onrender.com/)

## Overview

CodeBin is a backend API designed with modern engineering principles, featuring clean architecture, stateless JWT authentication, and enterprise-grade operational patterns. The service demonstrates professional software engineering practices including dependency injection, the repository pattern, structured logging, and comprehensive error handling.

## Key Technical Features

- **Clean Architecture**: Dependency injection pattern with clear separation of concerns between handlers, business logic, and data access layers
- **RESTful API Design**: Fully decoupled JSON API with expressive routing via `chi` router, including parameterized endpoints
- **Stateless Authentication**: JWT-based token authentication with bcrypt password hashing and custom middleware for protected routes
- **Repository Pattern**: Abstract data access layer for testability and maintainability
- **Performance Optimization**: In-memory time-based caching using `go-cache` for frequently accessed snippets, with optimized database indexes
- **Database Migrations**: Schema evolution managed through `soda` migrations for version-controlled database changes
- **Production-Ready Logging**: Structured JSON logging via Go's `slog` package for observability
- **High Availability**: Panic recovery middleware to prevent cascading failures
- **Externalized Configuration**: Server settings managed via command-line flags and environment variables
- **Comprehensive Testing**: Unit tests for API handlers using Go's standard `testing` and `httptest` packages

## Technology Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL 14+
- **Database Driver**: pgx/v5
- **Router**: go-chi/chi/v5
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt
- **Caching**: go-cache
- **Migrations**: gobuffalo/pop (soda CLI)
- **Logging**: slog (standard library)

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Soda CLI tool for migrations

### Installing Soda

```bash
go install github.com/gobuffalo/pop/v6/soda@latest
```

## Local Setup

### 1. Clone the Repository

```bash
git clone https://github.com/aakash-1857/codebin.git
cd codebin
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Database Setup

Create a PostgreSQL database and user:

```bash
# Connect to PostgreSQL as superuser
psql -U postgres

# In the PostgreSQL shell:
CREATE DATABASE codebin;
CREATE USER codebin_user WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE codebin TO codebin_user;

# For PostgreSQL 15+, also grant schema privileges:
\c codebin
GRANT ALL ON SCHEMA public TO codebin_user;
```

### 4. Configure Database Connection

Edit the `database.yml` file to match your local PostgreSQL configuration:

```yaml
development:
  dialect: postgres
  database: codebin
  user: codebin_user
  password: your_secure_password
  host: 127.0.0.1
  port: 5432
  pool: 5
```

### 5. Run Migrations

```bash
soda migrate up
```

This will create the necessary tables and indexes.

### 6. Run the Application

```bash
# Run with default settings (port 4000, development environment)
go run ./cmd/api

# Or with custom configuration:
go run ./cmd/api -port=8080 -env=production -db-dsn="postgres://codebin_user:your_password@localhost/codebin"
```

### 7. Verify the Server is Running

```bash
curl http://localhost:4000/v1/healthcheck
```

Expected response:
```json
{
  "status": "available",
  "environment": "development",
  "version": "1.0.0"
}
```

## API Endpoints

### Public Endpoints

- `GET /v1/healthcheck` - Health check endpoint
- `POST /v1/users` - Register a new user
- `POST /v1/tokens/authentication` - Authenticate and receive JWT

### Protected Endpoints (Require JWT)

- `POST /v1/snippets` - Create a new code snippet
- `GET /v1/snippets/{id}` - Retrieve a specific snippet (cached)
- `GET /v1/snippets` - List all snippets

## Configuration Options

The application accepts the following command-line flags:

- `-port`: HTTP server port (default: 4000)
- `-env`: Environment name (default: development)
- `-db-dsn`: PostgreSQL connection string

Example:
```bash
go run ./cmd/api -port=8080 -env=production -db-dsn="postgres://user:pass@localhost/codebin"
```



## Project Structure

```
│   codebin.exe
│   database.yml
│   go.mod
│   go.sum
│   payload.json
│   README.md
│
├───cmd
│   └───web
│           handlers.go
│           handlers_test.go
│           helpers.go
│           main.go
│           main_test.go
│           middleware.go
│           routes.go
│
├───internal
│   ├───models
│   │       snippet.go
│   │       user.go
│   │
│   └───repository
│           postgres.go
│           user_repository.go
│
└───migrations
        20250930132226_create_snippets_table.down.fizz
        20250930132226_create_snippets_table.up.fizz
        20250930134739_add_uuid_default_to_snippets.down.fizz
        20250930134739_add_uuid_default_to_snippets.up.fizz
        20250930143846_set_timestamp_defaults.down.fizz
        20250930143846_set_timestamp_defaults.up.fizz
        20251001050610_add_snippets_created_at_index.down.fizz
        20251001050610_add_snippets_created_at_index.up.fizz
        20251001051728_create_users_table.down.fizz
        20251001051728_create_users_table.up.fizz
        schema.sql
```

## Security Considerations

- Passwords are hashed using bcrypt with a cost factor of 12
- JWTs are signed using HMAC-SHA256
- Protected endpoints validate JWT signatures and expiration
- SQL injection protection via parameterized queries (pgx)
- Panic recovery middleware prevents service disruption

## Performance Features

- Database index on `snippets.created_at DESC` for optimized listing queries
- In-memory caching with TTL for frequently accessed snippets
- Connection pooling via pgx


## Author

Your Name - [your.email@example.com](mailto:akh.kashi.g@protonmail.com)

## Acknowledgments

Built with Go and inspired by modern backend engineering best practices.
