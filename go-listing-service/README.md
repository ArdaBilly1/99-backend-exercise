# Go Listing Service - Clean Architecture

A simple listing service API implemented in Go following Clean Architecture principles.

## Architecture

This project demonstrates Clean Architecture with clear separation of concerns:

```
Domain Layer (Core Business Logic)
    ↓
Repository Layer (Data Persistence Interface)
    ↓
Use Case Layer (Business Logic Orchestration)
    ↓
Handler Layer (HTTP/Transport)
```

### Layers

1. **Domain Layer** (`internal/domain/`)
   - Core entities and business rules
   - No external dependencies
   - Contains `Listing` entity with validation logic

2. **Repository Layer** (`internal/repository/`)
   - Data persistence abstraction
   - Interface definition + SQLite implementation
   - Handles database operations

3. **Use Case Layer** (`internal/usecase/`)
   - Business logic orchestration
   - Depends only on domain and repository interfaces
   - Coordinates between layers

4. **Handler Layer** (`internal/handler/`)
   - HTTP request/response handling
   - Form parsing and JSON serialization
   - Depends on use case interfaces

### Dependency Injection

All dependencies flow inward: Handler → Use Case → Repository → Domain

The `main.go` wires everything together:
```go
repo := sqlite.NewListingRepository()
useCase := usecase.NewListingUseCase(repo)
handler := handler.NewListingHandler(useCase)
```

## Project Structure

```
go-listing-service/
├── cmd/
│   └── server/
│       └── main.go              # Entry point with DI
├── internal/
│   ├── domain/
│   │   └── listing.go           # Core entity
│   ├── repository/
│   │   ├── listing_repository.go    # Interface
│   │   └── sqlite/
│   │       └── listing_repository.go # SQLite impl
│   ├── usecase/
│   │   └── listing_usecase.go       # Business logic
│   └── handler/
│       ├── response.go              # Response helpers
│       └── listing_handler.go       # HTTP handlers
├── go.mod
├── go.sum
└── README.md
```

## API Endpoints

### Create Listing
```bash
POST /listings
Content-Type: application/x-www-form-urlencoded

user_id=1&listing_type=rent&price=4500
```

**Response:**
```json
{
  "result": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "listing_type": "rent",
    "price": 4500,
    "created_at": 1234567890000000,
    "updated_at": 1234567890000000
  }
}
```

### Get Listings
```bash
GET /listings?page_num=1&page_size=10&user_id=1
```

**Response:**
```json
{
  "result": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "listing_type": "rent",
      "price": 4500,
      "created_at": 1234567890000000,
      "updated_at": 1234567890000000
    }
  ]
}
```

### Health Check
```bash
GET /listings/ping
```

**Response:**
```json
{
  "result": true,
  "data": {
    "status": "ok"
  }
}
```

## Getting Started

### Prerequisites
- Go 1.21 or higher
- SQLite3

### Installation

```bash
# Navigate to project directory
cd go-listing-service

# Install dependencies
go mod download
```

### Running the Service

```bash
# Run with default settings (port 6000, debug mode)
go run cmd/server/main.go

# Run with custom port
go run cmd/server/main.go -port=8888

# Run in production mode (disable debug)
go run cmd/server/main.go -debug=false

# Build and run binary
go build -o listing-service cmd/server/main.go
./listing-service -port=6000
```

### Testing

```bash
# Create a listing
curl localhost:6000/listings -XPOST \
    -d user_id=1 \
    -d listing_type=rent \
    -d price=4500

# Get all listings
curl "localhost:6000/listings?page_num=1&page_size=10"

# Get listings by user
curl "localhost:6000/listings?user_id=1"

# Health check
curl localhost:6000/listings/ping
```

## Validation Rules

- `listing_type`: Must be either "rent" or "sale"
- `price`: Must be greater than 0
- `user_id`: Required integer
- Timestamps are automatically generated in microseconds

## Database Schema

**Table: listings**
- `id` - INTEGER PRIMARY KEY AUTOINCREMENT
- `user_id` - INTEGER NOT NULL
- `listing_type` - TEXT NOT NULL
- `price` - INTEGER NOT NULL
- `created_at` - INTEGER NOT NULL (microseconds)
- `updated_at` - INTEGER NOT NULL (microseconds)

The database file `listings.db` is automatically created on first run.

## Benefits of This Architecture

1. **Testability**: Each layer can be tested independently with mocks
2. **Maintainability**: Clear separation of concerns
3. **Flexibility**: Easy to swap implementations (e.g., PostgreSQL instead of SQLite)
4. **Scalability**: Business logic is decoupled from frameworks
5. **Independence**: Domain logic doesn't depend on external packages
