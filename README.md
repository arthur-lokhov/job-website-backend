# Job Website Backend

This is the backend for a job website, built with Go. It follows a clean architecture, separating concerns into Handlers, Services, Repositories, and Models. The application leverages Docker and Docker Compose for containerization, PostgreSQL as its primary database, and Valkey (a Redis fork) for caching.

## Building and Running

To build and run the application, you need Docker and Docker Compose installed on your system.

All necessary configuration parameters are defined as environment variables within `docker-compose.yml` with sensible default values. An optional `.env` file can be used to override these defaults for local development.

1.  **Start the application:**

    ```bash
    docker-compose up -d
    ```

    This command will build the `job-website-backend:latest` image (if not already built) and start all defined services (database, cache, and the application).

2.  **Access the API:**

    The API is documented using Swagger. You can access the Swagger UI at `http://localhost:8080/swagger/index.html`.

    Key API Endpoints:
    *   `GET /api/v1/vacancies`: Get a list of active vacancies.
    *   `GET /api/v1/vacancies/{id}`: Get a vacancy by ID.
    *   `GET /api/v1/vacancies/{id}/form`: Get the application form for a vacancy.
    *   `POST /api/v1/applications`: Create a new application.
    *   `GET /api/v1/departments`: Get a list of departments.
    *   `GET /api/v1/levels`: Get a list of levels.
    *   `GET /api/v1/locations`: Get a list of locations.

    Admin Endpoints (require authentication):
    *   `GET /api/v1/admin/vacancies`: Get a list of all vacancies.
    *   `POST /api/v1/admin/vacancies`: Create a new vacancy.
    *   `PATCH /api/v1/admin/vacancies/{id}`: Update a vacancy.
    *   `DELETE /api/v1/admin/vacancies/{id}`: Delete a vacancy.
    *   `GET /api/v1/admin/applications`: Get a list of all applications.
    *   `PATCH /api/v1/admin/applications/{id}`: Update an application.
    *   `DELETE /api/v1/admin/applications/{id}`: Delete an application.

## Development Conventions

*   **Testing:** To run the Go tests, use the command:
    ```bash
    go test ./...
    ```
*   **Configuration:** Configuration is managed exclusively through environment variables. The application uses `viper` to load these variables, with explicit binding for robustness. Default values are provided in `docker-compose.yml`, which can be overridden by a local `.env` file.
*   **Database Migrations:** The project uses Goose for database migrations. Migrations are automatically run on application startup.
*   **UUID Support:** PostgreSQL UUID support is enabled via the `uuid-ossp` extension, which is configured in `scripts/init_dev.sql` for development environments.

## Authentication and Authorization

The application integrates with an external AuthService for authentication and authorization.

*   **Public Key Validation:** The service fetches the public key from the AuthService to validate JWT tokens. This key is periodically refreshed.
*   **Blacklist Check:** JWT tokens are checked against a blacklist maintained by the AuthService. The blacklist is also periodically refreshed.
*   **Permissions:** User permissions are fetched from the AuthService and cached for 24 hours to reduce load on the AuthService.
*   **API Token:** The service uses its own API token to communicate with the AuthService for administrative tasks (e.g., fetching permissions).

## Example cURL Commands

### Public Endpoints

**Get all vacancies:**
```bash
curl -X GET "http://localhost:8080/api/v1/vacancies"
```

**Get important vacancies:**
```bash
curl -X GET "http://localhost:8080/api/v1/vacancies?important=true"
```

**Search vacancies by name (e.g., "Engineer"):**
```bash
curl -X GET "http://localhost:8080/api/v1/vacancies?search=Engineer"
```

**Get a vacancy by ID (replace with actual ID):**
```bash
curl -X GET "http://localhost:8080/api/v1/vacancies/YOUR_VACANCY_ID"
```

**Get application form for a vacancy (replace with actual ID):**
```bash
curl -X GET "http://localhost:8080/api/v1/vacancies/YOUR_VACANCY_ID/form"
```

**Create a new application:**
```bash
curl -X POST "http://localhost:8080/api/v1/applications" \
     -H "Content-Type: application/json" \
     -d '{
           "vacancy_id": "YOUR_VACANCY_ID",
           "answers": {
             "fullName": "John Doe",
             "email": "john.doe@example.com",
             "phone": "123-456-7890"
           }
         }'
```

**Get all departments:**
```bash
curl -X GET "http://localhost:8080/api/v1/departments"
```

**Get all levels:**
```bash
curl -X GET "http://localhost:8080/api/v1/levels"
```

**Get all locations:**
```bash
curl -X GET "http://localhost:8080/api/v1/locations"
```

### Admin Endpoints (requires Authorization header with JWT token)

**Get all admin vacancies:**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/vacancies" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Get all admin applications:**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/applications" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Search admin applications by FIO (e.g., "John Doe"):**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/applications?search=John%20Doe" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
