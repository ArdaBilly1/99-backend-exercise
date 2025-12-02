# Microservices System Guide

This guide explains how to run the complete microservices architecture for the property rental/sales platform.

## System Architecture

```
┌─────────────────────────────────────────┐
│        External Clients                 │
│    (Mobile Apps, Websites)              │
└────────────────┬────────────────────────┘
                 │ JSON (port 8000)
                 ▼
┌─────────────────────────────────────────┐
│         go-public-api                   │
│         (API Gateway)                   │
│    - Aggregates service data            │
│    - Enriches responses                 │
│    - JSON requests/responses            │
└───────────┬──────────────┬──────────────┘
            │              │
   ┌────────┘              └────────┐
   │ form-encoded                   │ form-encoded
   │ (port 6000)                    │ (port 7000)
   ▼                                ▼
┌──────────────────-┐      ┌──────────────────┐
│ go-listing-service│      │  go-user-service │
│  (Internal API)   │      │  (Internal API)  │
│ - Listing CRUD    │      │  - User CRUD     │
│ - SQLite DB       │      │  - SQLite DB     │
└──────────────────-┘      └──────────────────┘
```

## Services Overview

### 1. Listing Service (Port 6000)
- **Purpose**: Manages property listings (rent/sale)
- **Database**: `listings.db` (SQLite)
- **API Format**: Form-encoded requests, JSON responses
- **Technology**: Python/Tornado
- **File**: `listing_service.py`

### 2. User Service (Port 7000)
- **Purpose**: Manages user data
- **Database**: `users.db` (SQLite)
- **API Format**: Form-encoded requests, JSON responses
- **Directory**: `go-user-service/`

### 3. Public API (Port 8000)
- **Purpose**: External-facing gateway that aggregates listing + user data
- **Database**: None (aggregates from other services)
- **API Format**: JSON requests and responses
- **Directory**: `go-public-api/`

## Running All Services

### Method 1: Using Docker (Recommended)

The easiest way to run all services is using Docker:

```bash
# Quick start (uses the helper script)
./start-services.sh

# Or manually with docker-compose
docker-compose up -d --build

# Check status
docker-compose ps

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

See [DOCKER_GUIDE.md](DOCKER_GUIDE.md) for complete Docker documentation.

### Method 2: Using Multiple Terminals (Local Development)

**Terminal 1 - Listing Service:**
```bash
# Activate Python virtual environment (if not already active)
source env/bin/activate

# Run listing service
python listing_service.py --port=6000 --debug=true
```

**Terminal 2 - User Service:**
```bash
cd go-user-service
go run cmd/server/main.go --port=7000 --debug=true
```

**Terminal 3 - Public API:**
```bash
cd go-public-api
go run cmd/server/main.go --port=8000 --debug=true
```

### Method 2: Using Built Binaries

```bash
# Build all services
cd go-listing-service && go build -o bin/listing-service cmd/server/main.go
cd ../go-user-service && go build -o bin/user-service cmd/server/main.go
cd ../go-public-api && go build -o bin/public-api cmd/server/main.go

# Run in separate terminals
./go-listing-service/bin/listing-service --port=6000 &
./go-user-service/bin/user-service --port=7000 &
./go-public-api/bin/public-api --port=8000 &
```

## Testing the Complete System

### Quick Health Check
```bash
# Check all services are running
curl localhost:6000/listings/ping  # Should return "pong!"
curl localhost:7000/users/ping     # Should return "pong!"
curl localhost:8000/public-api/ping # Should return "pong!"
```

### End-to-End Workflow

**1. Create a User via Public API:**
```bash
curl localhost:8000/public-api/users \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Johnson"}'
```

Response:
```json
{
  "user": {
    "id": 1,
    "name": "Alice Johnson",
    "created_at": 1733065800000000,
    "updated_at": 1733065800000000
  }
}
```

**2. Create a Listing for that User:**
```bash
curl localhost:8000/public-api/listings \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"listing_type":"rent","price":5500}'
```

Response:
```json
{
  "listing": {
    "id": 1,
    "user_id": 1,
    "listing_type": "rent",
    "price": 5500,
    "created_at": 1733065800000000,
    "updated_at": 1733065800000000
  }
}
```

**3. Get Enriched Listings (with embedded user data):**
```bash
curl localhost:8000/public-api/listings
```

Response:
```json
{
  "result": true,
  "listings": [
    {
      "id": 1,
      "listing_type": "rent",
      "price": 5500,
      "created_at": 1733065800000000,
      "updated_at": 1733065800000000,
      "user": {
        "id": 1,
        "name": "Alice Johnson",
        "created_at": 1733065800000000,
        "updated_at": 1733065800000000
      }
    }
  ]
}
```

### Running Automated Tests

Each service has a test script:

```bash
# Test listing service (internal API)
cd go-listing-service
./test-api.sh

# Test user service (internal API)
cd go-user-service
./test-api.sh

# Test public API (external API - requires other services running)
cd go-public-api
./test-api.sh
```

## API Comparison

### Internal Services (Form-Encoded)
Both listing and user services use form-encoded requests:

```bash
curl localhost:6000/listings -XPOST \
  -d user_id=1 \
  -d listing_type=rent \
  -d price=5500
```

### Public API (JSON)
The public API uses JSON for external clients:

```bash
curl localhost:8000/public-api/listings -XPOST \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"listing_type":"rent","price":5500}'
```

## Key Architectural Principles

1. **Service Independence**: Each service has its own database and can be deployed independently
2. **No Direct DB Access**: Services never access each other's databases directly
3. **API Gateway Pattern**: Public API aggregates data from multiple services
4. **Protocol Translation**: Public API converts between JSON (external) and form-encoded (internal)
5. **Data Enrichment**: Public API embeds user information into listing responses

## Port Summary

| Service        | Port | Protocol      | Purpose           |
|----------------|------|---------------|-------------------|
| Listing Service| 6000 | Form-encoded  | Internal API      |
| User Service   | 7000 | Form-encoded  | Internal API      |
| Public API     | 8000 | JSON          | External Gateway  |

## Troubleshooting

### Service Won't Start
- Check if port is already in use: `lsof -i :PORT_NUMBER`
- Verify Go is installed: `go version`
- Check dependencies: `go mod download`

### Public API Can't Connect to Services
- Ensure listing service is running on port 6000
- Ensure user service is running on port 7000
- Check service URLs in public API startup logs

### Database Issues
- Delete database files to reset: `rm go-listing-service/listings.db go-user-service/users.db`
- Databases are auto-created on first run

## Development Tips

### Using Different Ports
```bash
go run cmd/server/main.go --port=CUSTOM_PORT
```

### Disabling Debug Mode
```bash
go run cmd/server/main.go --debug=false
```

### Custom Service URLs for Public API
```bash
go run cmd/server/main.go \
  --listing-service=http://custom-host:6000 \
  --user-service=http://custom-host:7000
```

## Next Steps

- Add authentication/authorization
- Implement caching layer
- Add service discovery
- Implement circuit breakers for resilience
- Add distributed tracing
- Implement async communication with message queues

## Resources

- [Listing Service README](go-listing-service/README.md)
- [User Service README](go-user-service/README.md)
- [Public API README](go-public-api/README.md)
- [Main Project README](README.md)
