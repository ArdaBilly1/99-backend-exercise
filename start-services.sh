#!/bin/bash

echo "=================================="
echo "Starting Microservices System"
echo "=================================="
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Error: docker-compose is not installed. Please install it and try again."
    exit 1
fi

echo "✓ Docker is running"
echo "✓ docker-compose is available"
echo ""

# Stop any existing containers
echo "Stopping existing containers..."
docker-compose down

# Build and start services
echo ""
echo "Building and starting services..."
echo "(This may take a few minutes on first run)"
echo ""
docker-compose up -d --build

# Wait for services to be healthy
echo ""
echo "Waiting for services to be healthy..."
echo "(This may take 10-30 seconds)"
echo ""

MAX_WAIT=60
WAITED=0

while [ $WAITED -lt $MAX_WAIT ]; do
    LISTING_HEALTH=$(docker inspect --format='{{.State.Health.Status}}' listing-service 2>/dev/null || echo "starting")
    USER_HEALTH=$(docker inspect --format='{{.State.Health.Status}}' user-service 2>/dev/null || echo "starting")
    PUBLIC_HEALTH=$(docker inspect --format='{{.State.Health.Status}}' public-api 2>/dev/null || echo "starting")
    
    if [ "$LISTING_HEALTH" = "healthy" ] && [ "$USER_HEALTH" = "healthy" ] && [ "$PUBLIC_HEALTH" = "healthy" ]; then
        echo ""
        echo "✓ All services are healthy!"
        break
    fi
    
    echo -n "."
    sleep 2
    WAITED=$((WAITED + 2))
done

echo ""
echo ""

# Check final status
docker-compose ps

echo ""
echo "=================================="
echo "Service Status"
echo "=================================="
echo ""

# Test endpoints
echo "Testing service endpoints..."
echo ""

if curl -sf http://localhost:6000/listings/ping > /dev/null; then
    echo "✓ Listing Service (http://localhost:6000) - OK"
else
    echo "❌ Listing Service (http://localhost:6000) - FAILED"
fi

if curl -sf http://localhost:7000/users/ping > /dev/null; then
    echo "✓ User Service (http://localhost:7000) - OK"
else
    echo "❌ User Service (http://localhost:7000) - FAILED"
fi

if curl -sf http://localhost:8000/public-api/ping > /dev/null; then
    echo "✓ Public API (http://localhost:8000) - OK"
else
    echo "❌ Public API (http://localhost:8000) - FAILED"
fi

echo ""
echo "=================================="
echo "Quick Test Commands"
echo "=================================="
echo ""
echo "Create a user:"
echo "  curl localhost:8000/public-api/users -X POST -H 'Content-Type: application/json' -d '{\"name\":\"Alice\"}'"
echo ""
echo "Create a listing:"
echo "  curl localhost:8000/public-api/listings -X POST -H 'Content-Type: application/json' -d '{\"user_id\":1,\"listing_type\":\"rent\",\"price\":5500}'"
echo ""
echo "Get enriched listings:"
echo "  curl localhost:8000/public-api/listings | jq ."
echo ""
echo "=================================="
echo "Useful Commands"
echo "=================================="
echo ""
echo "View logs:"
echo "  docker-compose logs -f"
echo ""
echo "Stop services:"
echo "  docker-compose down"
echo ""
echo "Restart services:"
echo "  docker-compose restart"
echo ""
echo "See DOCKER_GUIDE.md for more commands"
echo ""
