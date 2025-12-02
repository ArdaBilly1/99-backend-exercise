#!/bin/bash

echo "==================================="
echo "Starting Microservices (Local Mode)"
echo "==================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed. Please install Go and try again."
    echo "   Visit: https://go.dev/doc/install"
    exit 1
fi

# Check if Python 3 is installed
if ! command -v python3 &> /dev/null; then
    echo "❌ Error: Python 3 is not installed. Please install Python 3 and try again."
    exit 1
fi

# Check if virtualenv exists
if [ ! -d "env" ]; then
    echo "⚠️  Virtual environment not found. Creating..."
    if ! command -v virtualenv &> /dev/null; then
        echo "❌ Error: virtualenv is not installed."
        echo "   Install with: pip3 install virtualenv"
        exit 1
    fi
    virtualenv env --python=$(which python3)
    source env/bin/activate
    pip install -r python-libs.txt
    deactivate
fi

echo "✓ Go is installed"
echo "✓ Python 3 is installed"
echo "✓ Virtual environment ready"
echo ""

# Create log directory
mkdir -p logs

# Kill any existing services on these ports
echo "Stopping any existing services on ports 6000, 7000, 8000..."
lsof -ti:6000 | xargs kill -9 2>/dev/null
lsof -ti:7000 | xargs kill -9 2>/dev/null
lsof -ti:8000 | xargs kill -9 2>/dev/null
sleep 1

echo ""
echo "Starting services..."
echo ""

# Start Listing Service (Python)
echo "▶ Starting Listing Service on port 6000..."
source env/bin/activate
nohup python3 listing_service.py --port=6000 --debug=true > logs/listing-service.log 2>&1 &
LISTING_PID=$!
deactivate
echo "  PID: $LISTING_PID"
sleep 2

# Start User Service (Go)
echo "▶ Starting User Service on port 7000..."
cd go-user-service
nohup go run cmd/server/main.go --port=7000 --debug=true > ../logs/user-service.log 2>&1 &
USER_PID=$!
cd ..
echo "  PID: $USER_PID"
sleep 2

# Start Public API (Go)
echo "▶ Starting Public API on port 8000..."
cd go-public-api
nohup go run cmd/server/main.go --port=8000 --debug=true > ../logs/public-api.log 2>&1 &
PUBLIC_PID=$!
cd ..
echo "  PID: $PUBLIC_PID"

echo ""
echo "Waiting for services to be ready..."
sleep 5

echo ""
echo "==================================="
echo "Service Status"
echo "==================================="
echo ""

# Test endpoints
LISTING_STATUS="❌ FAILED"
USER_STATUS="❌ FAILED"
PUBLIC_STATUS="❌ FAILED"

if curl -sf http://localhost:6000/listings/ping > /dev/null 2>&1; then
    LISTING_STATUS="✓ OK"
fi

if curl -sf http://localhost:7000/users/ping > /dev/null 2>&1; then
    USER_STATUS="✓ OK"
fi

if curl -sf http://localhost:8000/public-api/ping > /dev/null 2>&1; then
    PUBLIC_STATUS="✓ OK"
fi

echo "Listing Service (http://localhost:6000) - $LISTING_STATUS"
echo "User Service    (http://localhost:7000) - $USER_STATUS"
echo "Public API      (http://localhost:8000) - $PUBLIC_STATUS"

echo ""
echo "==================================="
echo "Process IDs"
echo "==================================="
echo ""
echo "Listing Service PID: $LISTING_PID"
echo "User Service PID:    $USER_PID"
echo "Public API PID:      $PUBLIC_PID"

# Save PIDs to file for easy stopping
echo "$LISTING_PID" > .pids
echo "$USER_PID" >> .pids
echo "$PUBLIC_PID" >> .pids

echo ""
echo "==================================="
echo "Quick Test Commands"
echo "==================================="
echo ""
echo "Create a user:"
echo "  curl localhost:8000/public-api/users -X POST -H 'Content-Type: application/json' -d '{\"name\":\"Alice\"}'"
echo ""
echo "Create a listing:"
echo "  curl localhost:8000/public-api/listings -X POST -H 'Content-Type: application/json' -d '{\"user_id\":1,\"listing_type\":\"rent\",\"price\":5500}'"
echo ""
echo "Get enriched listings:"
echo "  curl localhost:8000/public-api/listings"
echo ""
echo "==================================="
echo "Useful Commands"
echo "==================================="
echo ""
echo "View logs:"
echo "  tail -f logs/listing-service.log"
echo "  tail -f logs/user-service.log"
echo "  tail -f logs/public-api.log"
echo ""
echo "Stop all services:"
echo "  ./stop-local.sh"
echo ""
