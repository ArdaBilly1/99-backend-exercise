# Quick Start Guide

## Running the Service

### Option 1: Run directly with Go
```bash
cd go-user-service
go run cmd/server/main.go
```

The service will start on port **7000** by default.

### Option 2: Run with custom port
```bash
go run cmd/server/main.go --port=8000 --debug=false
```

### Option 3: Use the compiled binary
```bash
./bin/user-service --port=7000
```

## Testing the API

### Manual Testing with curl

**Health Check:**
```bash
curl localhost:7000/users/ping
```

**Create a User:**
```bash
curl localhost:7000/users -XPOST -d name="John Doe"
```

**Get Specific User:**
```bash
curl localhost:7000/users/1
```

**Get All Users:**
```bash
curl "localhost:7000/users?page_num=1&page_size=10"
```

### Automated Testing

Run the provided test script (requires `jq` for JSON formatting):
```bash
./test-api.sh
```

## Building the Service

**Build for current platform:**
```bash
go build -o bin/user-service cmd/server/main.go
```

**Build with all dependencies:**
```bash
go mod download
go build -o bin/user-service cmd/server/main.go
```

## Database

The service automatically creates a `users.db` SQLite database file in the current directory on first run. The database schema is initialized automatically.

## Default Configuration

- **Port**: 7000
- **Debug Mode**: true
- **Database**: users.db (SQLite)
- **Request Format**: application/x-www-form-urlencoded
- **Response Format**: application/json

## Example Response Format

**Success:**
```json
{
    "result": true,
    "data": {
        "user": {
            "id": 1,
            "name": "John Doe",
            "created_at": 1733065800000000,
            "updated_at": 1733065800000000
        }
    }
}
```

**Error:**
```json
{
    "result": false,
    "errors": ["name is required"]
}
```

## Stopping the Service

Press `Ctrl+C` to gracefully shutdown the service.
