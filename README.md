# Job Website Backend

A Go-based backend service for managing job vacancies and applications.

## Features

- RESTful API for managing job vacancies and applications
- JWT-based authentication with public key validation
- Role-based access control
- PostgreSQL database with migrations
- CORS support
- Graceful shutdown
- Configuration management with Viper

## Prerequisites

- Go 1.21 or later
- PostgreSQL 12 or later
- Make (optional, for using Makefile commands)

## Project Structure

```
.
├── cmd/
│   └── api/              # Application entry point
├── config/              # Configuration files
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models
│   ├── repositories/    # Database repositories
│   ├── router/          # HTTP router
│   └── services/        # Business logic
├── migrations/          # Database migrations
├── config.yaml         # Default configuration
├── go.mod             # Go module file
├── go.sum             # Go module checksum
└── README.md          # This file
```

## Setup

1. Clone the repository:
   ```bash
   git clone https://git.cyberzone.dev/project-vacancy-website/job-website-backend.git
   cd job-website-backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Create a PostgreSQL database:
   ```bash
   createdb job_website
   ```

4. Run database migrations:
   ```bash
   go run cmd/migrate/main.go
   ```

5. Copy the default configuration and modify as needed:
   ```bash
   cp config/config.yaml.example config/config.yaml
   ```

6. Place your JWT public key in `config/public.pem`

## Running the Application

```bash
go run cmd/api/main.go
```

The server will start on port 8080 by default. You can change this in the configuration file.

## API Endpoints

### Public Endpoints

- `GET /api/v1/vacancies` - List all active vacancies
- `GET /api/v1/vacancies/{id}` - Get vacancy details
- `GET /api/v1/vacancies/{id}/form` - Get vacancy application form
- `POST /api/v1/applications` - Submit a job application

### Admin Endpoints (Requires Authentication)

- `GET /api/v1/admin/vacancies` - List all vacancies
- `POST /api/v1/admin/vacancies` - Create a new vacancy
- `PUT /api/v1/admin/vacancies/{id}` - Update a vacancy
- `DELETE /api/v1/admin/vacancies/{id}` - Delete a vacancy
- `GET /api/v1/admin/applications` - List all applications
- `PUT /api/v1/admin/applications/{id}` - Update application status
- `DELETE /api/v1/admin/applications/{id}` - Delete an application

## Authentication

The API uses JWT tokens for authentication. Admin endpoints require a valid JWT token with the "admin" permission in the Authorization header:

```
Authorization: Bearer <token>
```

## Configuration

The application can be configured using the `config.yaml` file or environment variables. The following settings are available:

- `server.port` - HTTP server port
- `server.readTimeout` - Read timeout in seconds
- `server.writeTimeout` - Write timeout in seconds
- `server.idleTimeout` - Idle timeout in seconds
- `database.url` - PostgreSQL connection URL
- `jwt.publicKeyPath` - Path to JWT public key
- `jwt.blacklistTTL` - Token blacklist TTL in hours
- `cors.*` - CORS configuration

## Development

### Running Tests

```bash
go test ./...
```

### Running Linter

```bash
golangci-lint run
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
