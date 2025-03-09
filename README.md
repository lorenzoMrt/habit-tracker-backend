# Habit Tracker Backend

A simple RESTful API backend for tracking habits, built with Go and PostgreSQL.

## Overview

This project provides a backend service for a habit tracking application. It allows users to:

- Create new habits
- List all habits
- Mark habits as completed

The application is containerized using Docker and can be easily deployed using Docker Compose.

## Tech Stack

- **Go**: Backend language (v1.22)
- **PostgreSQL**: Database
- **Gorilla Mux**: HTTP router
- **Docker & Docker Compose**: Containerization and orchestration

## Project Structure

```
habit-tracker-backend/
├── cmd/
│   └── api/
│       ├── main.go         # Main application entry point
│       └── main_test.go    # Tests for the API
├── infrastructure/
│   └── db/
│       └── migrations/     # Database migrations
│           └── migrations.sql
├── Dockerfile              # Docker image definition
├── compose.yml             # Docker Compose configuration
├── run-migrations.sh       # Script to run database migrations
├── go.mod                  # Go module definition
└── go.sum                  # Go module checksums
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/habits` | Create a new habit |
| GET | `/habits` | List all habits |
| PUT | `/habits/{id}/complete` | Mark a habit as completed |

## Prerequisites

- Docker and Docker Compose
- Go 1.23.4 (for local development)

## Getting Started

### Running with Docker Compose

1. Clone the repository:
   ```bash
   git clone https://github.com/lorenzoMrt/habit-tracker-backend.git
   cd habit-tracker-backend
   ```

2. Start the application:
   ```bash
   docker compose up
   ```

   This will:
   - Build the Go application
   - Start a PostgreSQL database
   - Run database migrations
   - Start the API server on port 8080

3. The API will be available at `http://localhost:8080`

### Local Development

1. Install Go 1.23.4
2. Install PostgreSQL
3. Create a database named `habit_tracker`
4. Run the migrations from `infrastructure/db/migrations/migrations.sql`
5. Update the database connection string in `cmd/api/main.go` if needed
6. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

## Example API Usage

### Create a new habit

```bash
curl -X POST http://localhost:8080/habits \
  -H "Content-Type: application/json" \
  -d '{"name": "Daily Exercise", "description": "30 minutes of exercise every day"}'
```

### List all habits

```bash
curl -X GET http://localhost:8080/habits
```

### Mark a habit as completed

```bash
curl -X PUT http://localhost:8080/habits/1/complete
```

## CORS Configuration

The API supports Cross-Origin Resource Sharing (CORS) to allow requests from different origins. By default, it allows requests from `http://localhost:5173`.

### Environment Variables

You can configure CORS settings using the following environment variables:

- `ALLOWED_ORIGINS`: Comma-separated list of allowed origins (e.g., `http://localhost:5173,http://localhost:3000`). Set to `*` to allow all origins.
- `PORT`: The port on which the server will listen (default: `8080`).
- `DATABASE_URL`: The PostgreSQL connection string. Make sure to include `sslmode=disable` if your PostgreSQL server doesn't have SSL enabled.

### Docker Compose Configuration

In the `compose.yml` file, you can configure these environment variables:

```yaml
services:
  app:
    # ...
    environment:
      - DATABASE_URL=postgres://postgres:password@db:5432/habit_tracker?sslmode=disable
      - ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000,http://frontend:5173
      - PORT=8080
```

## License

This project is licensed under the terms found in the LICENSE file.