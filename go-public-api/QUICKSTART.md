# Quick Start Guide - Public API

## Prerequisites

Make sure both internal services are running first:

```bash
# Terminal 1: Listing Service
cd go-listing-service
go run cmd/server/main.go --port=6000

# Terminal 2: User Service  
cd go-user-service
go run cmd/server/main.go --port=7000
```

## Running the Public API

### Option 1: Run directly with Go
```bash
cd go-public-api
go run cmd/server/main.go
```

The service will start on port **8000** by default.

### Option 2: Run with custom configuration
```bash
go run cmd/server/main.go \
  --port=9000 \
  --listing-service=http://localhost:6000 \
  --user-service=http://localhost:7000
```

### Option 3: Use the compiled binary
```bash
./bin/public-api --port=8000
```

## Quick API Tests

### 1. Health Check
```bash
curl localhost:8000/public-api/ping
```

### 2. Create a User
```bash
curl localhost:8000/public-api/users \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Smith"}'
```

### 3. Create a Listing
```bash
curl localhost:8000/public-api/listings \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"listing_type":"rent","price":5500}'
```

### 4. Get Enriched Listings
```bash
curl "localhost:8000/public-api/listings?page_num=1&page_size=10"
```

Notice how the response includes full user details embedded in each listing!

## Complete End-to-End Test

```bash
# 1. Create a user via public API
curl localhost:8000/public-api/users \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob Jones"}'
# Returns: {"user":{"id":1,"name":"Bob Jones",...}}

# 2. Create a listing for that user
curl localhost:8000/public-api/listings \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"listing_type":"sale","price":450000}'
# Returns: {"listing":{"id":1,"user_id":1,...}}

# 3. Get all listings (enriched with user data)
curl "localhost:8000/public-api/listings"
# Returns listings with embedded user objects
```

## Building the Service

```bash
go mod download
go build -o bin/public-api cmd/server/main.go
```

## Configuration

**Command-line Flags:**
- `--port` - Server port (default: 8000)
- `--debug` - Debug mode (default: true)
- `--listing-service` - Listing service URL (default: http://localhost:6000)
- `--user-service` - User service URL (default: http://localhost:7000)

## Response Formats

**Success Response:**
```json
{
    "result": true,
    "listings": [...]
}
```

**Error Response:**
```json
{
    "error": "error message"
}
```

## Stopping the Service

Press `Ctrl+C` to gracefully shutdown.
