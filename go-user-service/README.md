# Go User Service

A microservice for managing user data, built with Go using Clean Architecture principles.

## Architecture

This service follows Clean Architecture with clear separation of concerns:

- **Domain Layer** (`internal/domain/`) - Core business entities and validation rules
- **Repository Layer** (`internal/repository/`) - Data persistence interfaces and implementations
- **Use Case Layer** (`internal/usecase/`) - Business logic orchestration
- **Handler Layer** (`internal/handler/`) - HTTP request/response handling

## Technology Stack

- **Language**: Go 1.21+
- **HTTP Framework**: net/http with gorilla/mux for routing
- **Database**: SQLite3
- **Architecture**: Clean Architecture with dependency injection

## Project Structure

```
go-user-service/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/
│   │   └── user.go              # User entity and business rules
│   ├── repository/
│   │   ├── user_repository.go   # Repository interface
│   │   └── sqlite/
│   │       └── user_repository.go # SQLite implementation
│   ├── usecase/
│   │   └── user_usecase.go      # Business logic
│   └── handler/
│       ├── user_handler.go      # HTTP handlers
│       ├── response.go          # Response helpers
│       └── route.go             # Route configuration
├── go.mod
├── go.sum
└── README.md
```

## Database Schema

**Users Table** (SQLite in `users.db`):
- `id` (INTEGER, PRIMARY KEY, AUTOINCREMENT)
- `name` (TEXT, NOT NULL)
- `created_at` (INTEGER, NOT NULL) - Microseconds timestamp
- `updated_at` (INTEGER, NOT NULL) - Microseconds timestamp

## Getting Started

### Prerequisites

- Go 1.21 or higher
- SQLite3 (included with most systems)

### Installation

1. Navigate to the project directory:
```bash
cd go-user-service
```

2. Download dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o bin/user-service cmd/server/main.go
```

### Running the Service

Run with default settings (port 7000, debug mode):
```bash
go run cmd/server/main.go
```

Run with custom port:
```bash
go run cmd/server/main.go --port=8000 --debug=false
```

Or run the built binary:
```bash
./bin/user-service --port=7000 --debug=true
```

## API Endpoints

### Health Check
```bash
GET /users/ping
```

Response:
```
pong!
```

### Create User
```bash
POST /users
Content-Type: application/x-www-form-urlencoded

name=John Doe
```

Response:
```json
{
    "result": true,
    "data": {
        "user": {
            "id": 1,
            "name": "John Doe",
            "created_at": 1475820997000000,
            "updated_at": 1475820997000000
        }
    }
}
```

Example with curl:
```bash
curl localhost:7000/users -XPOST -d name="John Doe"
```

### Get Specific User
```bash
GET /users/{id}
```

Response:
```json
{
    "result": true,
    "data": {
        "user": {
            "id": 1,
            "name": "John Doe",
            "created_at": 1475820997000000,
            "updated_at": 1475820997000000
        }
    }
}
```

Example with curl:
```bash
curl localhost:7000/users/1
```

### Get All Users
```bash
GET /users?page_num=1&page_size=10
```

Query Parameters:
- `page_num` (int, default: 1) - Page number
- `page_size` (int, default: 10) - Items per page

Response:
```json
{
    "result": true,
    "data": {
        "users": [
            {
                "id": 1,
                "name": "John Doe",
                "created_at": 1475820997000000,
                "updated_at": 1475820997000000
            }
        ]
    }
}
```

Example with curl:
```bash
curl "localhost:7000/users?page_num=1&page_size=10"
```

## Error Handling

All errors follow a consistent format:
```json
{
    "result": false,
    "errors": ["error message here"]
}
```

## Development

### Running Tests
```bash
go test ./...
```

### Code Formatting
```bash
go fmt ./...
```

### Building for Production
```bash
CGO_ENABLED=1 go build -o bin/user-service cmd/server/main.go
```

Note: `CGO_ENABLED=1` is required for SQLite3 support.

## Design Principles

- **Dependency Injection**: All layers depend on interfaces, not concrete implementations
- **Single Responsibility**: Each layer has a clear, focused purpose
- **Testability**: Interfaces enable easy mocking for unit tests
- **Form-Encoded Requests**: Internal service uses `application/x-www-form-urlencoded` for consistency with other microservices
- **JSON Responses**: All responses are in JSON format with consistent structure

## License

This is an educational project for demonstrating microservices architecture.
