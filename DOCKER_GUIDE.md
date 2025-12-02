# Docker Deployment Guide

This guide explains how to run the complete microservices system using Docker and Docker Compose.

## Prerequisites

- Docker installed (version 20.10 or higher)
- Docker Compose installed (version 2.0 or higher)

Verify installation:
```bash
docker --version
docker-compose --version
```

## Architecture

The system consists of three containerized services:

```
┌──────────────────────────────────────────┐
│   External Clients (Port 8000)           │
└──────────────┬───────────────────────────┘
               │
               ▼
        ┌─────────────┐
        │ public-api  │ (Go, Port 8000)
        │  Container  │
        └──────┬──────┘
               │
        ┌──────┴──────┐
        │             │
   ┌────▼────┐   ┌───▼─────┐
   │ listing │   │  user   │
   │ service │   │ service │
   │(Python) │   │  (Go)   │
   │Port 6000│   │Port 7000│
   └─────────┘   └─────────┘
```

## Services

### 1. listing-service (Python/Tornado)
- **Image**: Built from `Dockerfile.listing`
- **Port**: 6000
- **Database**: SQLite (persisted in volume)
- **Base Image**: python:3.9-slim

### 2. user-service (Go)
- **Image**: Built from `Dockerfile.user`
- **Port**: 7000
- **Database**: SQLite (persisted in volume)
- **Base Image**: golang:1.21-alpine

### 3. public-api (Go)
- **Image**: Built from `Dockerfile.public`
- **Port**: 8000
- **Dependencies**: Waits for listing and user services to be healthy
- **Base Image**: golang:1.21-alpine

## Quick Start

### Build and Run All Services

```bash
# Build and start all services in detached mode
docker-compose up -d --build
```

### Check Service Status

```bash
# View all running containers
docker-compose ps

# View logs from all services
docker-compose logs

# View logs from a specific service
docker-compose logs listing-service
docker-compose logs user-service
docker-compose logs public-api

# Follow logs in real-time
docker-compose logs -f public-api
```

### Health Check

Wait for all services to be healthy (may take 10-30 seconds):

```bash
# Check health status
docker-compose ps

# Or test endpoints directly
curl localhost:6000/listings/ping  # Should return "pong!"
curl localhost:7000/users/ping     # Should return "pong!"
curl localhost:8000/public-api/ping # Should return "pong!"
```

## Testing the System

### Create a User

```bash
curl localhost:8000/public-api/users \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Johnson"}'
```

### Create a Listing

```bash
curl localhost:8000/public-api/listings \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"listing_type":"rent","price":5500}'
```

### Get Enriched Listings

```bash
curl localhost:8000/public-api/listings | jq .
```

## Docker Commands

### Stop All Services

```bash
docker-compose stop
```

### Start Existing Containers

```bash
docker-compose start
```

### Restart Services

```bash
# Restart all services
docker-compose restart

# Restart specific service
docker-compose restart public-api
```

### Stop and Remove Containers

```bash
docker-compose down
```

### Stop and Remove Everything (including volumes)

```bash
# WARNING: This will delete all data in databases
docker-compose down -v
```

### Rebuild Services

```bash
# Rebuild all services
docker-compose build

# Rebuild specific service
docker-compose build user-service

# Rebuild and restart
docker-compose up -d --build
```

## Viewing Logs

```bash
# All services
docker-compose logs

# Specific service
docker-compose logs listing-service
docker-compose logs user-service
docker-compose logs public-api

# Follow logs (real-time)
docker-compose logs -f

# Last 100 lines
docker-compose logs --tail=100

# Logs with timestamps
docker-compose logs -t
```

## Debugging

### Execute Commands in Running Containers

```bash
# Access listing service container
docker-compose exec listing-service sh

# Access user service container
docker-compose exec user-service sh

# Access public API container
docker-compose exec public-api sh
```

### View Container Details

```bash
# Inspect listing service
docker inspect listing-service

# View resource usage
docker stats
```

### Check Network Connectivity

```bash
# From public-api, test connection to listing-service
docker-compose exec public-api wget -O- http://listing-service:6000/listings/ping

# From public-api, test connection to user-service
docker-compose exec public-api wget -O- http://user-service:7000/users/ping
```

## Data Persistence

Database files are stored in Docker volumes:

```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect 99-backend-exercise_listing-data
docker volume inspect 99-backend-exercise_user-data

# Backup volumes
docker run --rm \
  -v 99-backend-exercise_listing-data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/listing-backup.tar.gz /data

# Restore volumes
docker run --rm \
  -v 99-backend-exercise_listing-data:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/listing-backup.tar.gz -C /
```

## Scaling Considerations

### Individual Service Scaling

```bash
# Scale user service to 3 instances
docker-compose up -d --scale user-service=3

# Note: You'll need a load balancer for this to work properly
```

### Environment Variables

You can customize service URLs via environment variables:

```bash
# In docker-compose.yml, public-api service:
environment:
  - LISTING_SERVICE_URL=http://listing-service:6000
  - USER_SERVICE_URL=http://user-service:7000
```

## Production Considerations

### Using Pre-built Images

For production, build and push images to a registry:

```bash
# Tag images
docker tag 99-backend-exercise_listing-service myregistry/listing-service:v1.0
docker tag 99-backend-exercise_user-service myregistry/user-service:v1.0
docker tag 99-backend-exercise_public-api myregistry/public-api:v1.0

# Push to registry
docker push myregistry/listing-service:v1.0
docker push myregistry/user-service:v1.0
docker push myregistry/public-api:v1.0
```

### Health Checks

All services have health checks configured:
- Listing Service: Checks `/listings/ping` every 10s
- User Service: Checks `/users/ping` every 10s
- Public API: Checks `/public-api/ping` every 10s

### Resource Limits

Add resource limits in docker-compose.yml:

```yaml
services:
  listing-service:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port
lsof -i :6000
lsof -i :7000
lsof -i :8000

# Kill process or use different ports in docker-compose.yml
```

### Service Won't Start

```bash
# Check logs
docker-compose logs service-name

# Check if container is running
docker-compose ps

# Restart specific service
docker-compose restart service-name
```

### Can't Connect Between Services

```bash
# Check network
docker network ls
docker network inspect 99-backend-exercise_microservices-network

# Verify service names resolve
docker-compose exec public-api ping listing-service
docker-compose exec public-api ping user-service
```

### Database Issues

```bash
# Remove volumes and start fresh
docker-compose down -v
docker-compose up -d --build
```

## Cleanup

### Remove Everything

```bash
# Stop and remove containers, networks, images, and volumes
docker-compose down -v --rmi all

# Remove dangling images
docker image prune -f

# Remove all unused data
docker system prune -a --volumes
```

## Development Workflow

### Hot Reload Development

For development with hot reload, mount source code as volumes:

```yaml
# Add to docker-compose.yml under listing-service
volumes:
  - ./listing_service.py:/app/listing_service.py
```

Then restart the service after code changes:
```bash
docker-compose restart listing-service
```

### Running Tests in Containers

```bash
# Run tests inside a service
docker-compose exec user-service go test ./...
```

## CI/CD Integration

### Example GitHub Actions Workflow

```yaml
name: Build and Test

on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build services
        run: docker-compose build
      - name: Start services
        run: docker-compose up -d
      - name: Wait for services
        run: sleep 30
      - name: Test services
        run: |
          curl -f http://localhost:8000/public-api/ping
      - name: Cleanup
        run: docker-compose down -v
```

## Summary

- **Build & Run**: `docker-compose up -d --build`
- **View Logs**: `docker-compose logs -f`
- **Stop**: `docker-compose down`
- **Clean Everything**: `docker-compose down -v --rmi all`

All services will be available at:
- Listing Service: http://localhost:6000
- User Service: http://localhost:7000
- Public API: http://localhost:8000
