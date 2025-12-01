# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview
This is a microservices architecture exercise for a property rental/sales platform. The system consists of three independent services that communicate via REST APIs:

1. **Listing Service** (`listing_service.py`) - Manages property listings (rent/sale)
2. **User Service** (not yet implemented) - Manages user data
3. **Public API Layer** (not yet implemented) - Aggregates data from listing and user services for external clients

**Key Architectural Principle**: Services are the sole gatekeepers to their databases. No direct database access is allowed between services—all communication must go through REST APIs.

## Development Commands

### Environment Setup
```bash
# Locate Python 3 installation
which python3

# Create virtual environment
virtualenv env --python=<path_to_python_3>

# Activate virtual environment
source env/bin/activate

# Install dependencies
pip install -r python-libs.txt

# Verify installation
pip freeze
```

### Running Services
```bash
# Run listing service (default port 6000, debug mode enabled)
python listing_service.py --port=6000 --debug=true

# Run with custom port
python listing_service.py --port=8888 --debug=false
```

### Testing Endpoints Manually
```bash
# Create a listing
curl localhost:6000/listings -XPOST \
    -d user_id=1 \
    -d listing_type=rent \
    -d price=4500

# Get all listings (with pagination)
curl "localhost:6000/listings?page_num=1&page_size=10"

# Get listings by user
curl "localhost:6000/listings?user_id=1"

# Health check
curl localhost:6000/listings/ping
```

## Architecture Details

### Service Communication Pattern
- **Listing Service & User Service**: Internal microservices that wrap their SQLite databases, exposing REST APIs with form-encoded requests (`application/x-www-form-urlencoded`)
- **Public API Layer**: External-facing gateway that aggregates data from internal services and returns JSON responses. It accepts JSON requests (`application/json`) and enriches listing data with user information by calling both services

### Data Flow Example
When a client requests listings via `/public-api/listings`:
1. Public API calls Listing Service's `GET /listings`
2. For each listing, Public API calls User Service's `GET /users/{user_id}`
3. Public API merges the data and returns enriched listing objects with embedded user details

### Database Schema
**Listings Table** (SQLite in `listings.db`):
- `id` (INTEGER, PRIMARY KEY, AUTOINCREMENT)
- `user_id` (INTEGER, NOT NULL)
- `listing_type` (TEXT, NOT NULL) - Must be "rent" or "sale"
- `price` (INTEGER, NOT NULL) - Must be > 0
- `created_at` (INTEGER, NOT NULL) - Microseconds timestamp
- `updated_at` (INTEGER, NOT NULL) - Microseconds timestamp

### Framework & Patterns
- **Framework**: Tornado (async web framework)
- **Database**: SQLite with `sqlite3.Row` factory for dict-like row access
- **Validation**: Input validation happens in handler methods before database operations
- **Timestamps**: All timestamps use microseconds (`int(time.time() * 1e6)`)
- **Pagination**: Implemented via `page_num` and `page_size` query parameters with `LIMIT`/`OFFSET`

### Key Implementation Notes
- All services use Tornado's `@tornado.gen.coroutine` decorator for async handlers
- Responses follow a consistent format: `{"result": true/false, "data": {...}}` or `{"result": false, "errors": [...]}`
- When implementing new services, maintain the established error handling pattern: collect validation errors in a list, return 400 status with errors array if validation fails
- The BaseHandler's `write_json()` method should be used for all JSON responses
