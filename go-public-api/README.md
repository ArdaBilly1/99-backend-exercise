# Go Public API

The external-facing API gateway for the property rental/sales platform. This service aggregates data from the listing and user services to provide enriched responses to external clients.

## Architecture

This service acts as an API Gateway that:
- Accepts JSON requests from external clients (mobile apps, websites)
- Communicates with internal microservices via their REST APIs
- Aggregates and enriches data before returning to clients
- Returns JSON responses in a format suitable for external consumption

### Key Design Principles

- **No Direct Database Access**: All data access goes through the listing and user service APIs
- **Service Aggregation**: Combines data from multiple services into single responses
- **Protocol Translation**: Accepts JSON from clients, communicates with services using form-encoded requests
- **Data Enrichment**: Embeds user information within listing responses

## Technology Stack

- **Language**: Go 1.21+
- **HTTP Framework**: net/http with gorilla/mux for routing
- **Architecture**: API Gateway pattern with HTTP clients for service communication

## Project Structure

```
go-public-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── model/
│   │   └── models.go            # Data models
│   ├── client/
│   │   ├── listing_client.go    # Listing service HTTP client
│   │   └── user_client.go       # User service HTTP client
│   └── handler/
│       ├── public_handler.go    # HTTP handlers
│       ├── response.go          # Response helpers
│       └── route.go             # Route configuration
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

The public API requires both internal services to be running:

1. **Listing Service** - Default: `http://localhost:6000`
2. **User Service** - Default: `http://localhost:7000`

## Getting Started

### Installation

1. Navigate to the project directory:
```bash
cd go-public-api
```

2. Download dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o bin/public-api cmd/server/main.go
```

### Running the Service

Run with default settings (port 8000):
```bash
go run cmd/server/main.go
```

Run with custom configuration:
```bash
go run cmd/server/main.go \
  --port=9000 \
  --listing-service=http://localhost:6000 \
  --user-service=http://localhost:7000 \
  --debug=true
```

Or run the built binary:
```bash
./bin/public-api --port=8000
```

### Command-line Flags

- `--port` (int, default: 8000) - Server port
- `--debug` (bool, default: true) - Debug mode
- `--listing-service` (string, default: "http://localhost:6000") - Listing service URL
- `--user-service` (string, default: "http://localhost:7000") - User service URL

## API Endpoints

### Health Check
```bash
GET /public-api/ping
```

Response:
```
pong!
```

### Get Listings (with enriched user data)
```bash
GET /public-api/listings?page_num=1&page_size=10&user_id=1
```

Query Parameters:
- `page_num` (int, default: 1) - Page number
- `page_size` (int, default: 10) - Items per page
- `user_id` (int, optional) - Filter by user ID

Response:
```json
{
    "result": true,
    "listings": [
        {
            "id": 1,
            "listing_type": "rent",
            "price": 6000,
            "created_at": 1475820997000000,
            "updated_at": 1475820997000000,
            "user": {
                "id": 1,
                "name": "John Doe",
                "created_at": 1475820997000000,
                "updated_at": 1475820997000000
            }
        }
    ]
}
```

Example with curl:
```bash
curl "localhost:8000/public-api/listings?page_num=1&page_size=10"
```

### Create User
```bash
POST /public-api/users
Content-Type: application/json
```

Request Body:
```json
{
    "name": "John Doe"
}
```

Response:
```json
{
    "user": {
        "id": 1,
        "name": "John Doe",
        "created_at": 1475820997000000,
        "updated_at": 1475820997000000
    }
}
```

Example with curl:
```bash
curl localhost:8000/public-api/users \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe"}'
```

### Create Listing
```bash
POST /public-api/listings
Content-Type: application/json
```

Request Body:
```json
{
    "user_id": 1,
    "listing_type": "rent",
    "price": 6000
}
```

Response:
```json
{
    "listing": {
        "id": 1,
        "user_id": 1,
        "listing_type": "rent",
        "price": 6000,
        "created_at": 1475820997000000,
        "updated_at": 1475820997000000
    }
}
```

Example with curl:
```bash
curl localhost:8000/public-api/listings \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"listing_type":"rent","price":6000}'
```

## Error Handling

All errors return a JSON response with an error message:
```json
{
    "error": "error message here"
}
```

## Data Flow Example

When a client requests `/public-api/listings`:

1. **Public API** receives JSON request from client
2. **Public API** calls Listing Service's `GET /listings` (form-encoded)
3. **Listing Service** returns listings with `user_id` fields
4. **Public API** calls User Service's `GET /users/{id}` for each unique user
5. **User Service** returns user details
6. **Public API** merges user data into listings
7. **Public API** returns enriched JSON response to client

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
go build -o bin/public-api cmd/server/main.go
```

## Running the Complete System

To run the entire microservices system:

```bash
# Terminal 1: Start listing service
cd go-listing-service
go run cmd/server/main.go --port=6000

# Terminal 2: Start user service
cd go-user-service
go run cmd/server/main.go --port=7000

# Terminal 3: Start public API
cd go-public-api
go run cmd/server/main.go --port=8000
```

## Design Decisions

- **JSON for External API**: Public API accepts/returns JSON for better client compatibility
- **Form-encoded for Internal**: Service-to-service communication uses form-encoded to match internal service contracts
- **Synchronous Aggregation**: User data is fetched synchronously for each listing (could be optimized with batch requests or caching)
- **Error Propagation**: Errors from internal services are propagated to clients with appropriate HTTP status codes

## License

This is an educational project for demonstrating microservices architecture.
