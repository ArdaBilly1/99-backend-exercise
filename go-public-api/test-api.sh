#!/bin/bash

# Test script for Public API
# Prerequisites: 
# - Listing service running on port 6000
# - User service running on port 7000
# - Public API running on port 8000

BASE_URL="http://localhost:8000"

echo "=================================="
echo "Testing Public API"
echo "=================================="
echo ""

# Test 1: Health check
echo "1. Health Check (GET /public-api/ping)"
curl -s ${BASE_URL}/public-api/ping
echo -e "\n"

# Test 2: Create users via public API
echo "2. Create Users (POST /public-api/users)"
echo "Creating user: John Doe"
USER1=$(curl -s ${BASE_URL}/public-api/users \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe"}')
echo $USER1 | jq .
USER1_ID=$(echo $USER1 | jq -r '.user.id')
echo ""

echo "Creating user: Jane Smith"
USER2=$(curl -s ${BASE_URL}/public-api/users \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Smith"}')
echo $USER2 | jq .
USER2_ID=$(echo $USER2 | jq -r '.user.id')
echo ""

# Test 3: Create listings via public API
echo "3. Create Listings (POST /public-api/listings)"
echo "Creating listing for user $USER1_ID"
curl -s ${BASE_URL}/public-api/listings \
  -X POST \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":$USER1_ID,\"listing_type\":\"rent\",\"price\":4500}" | jq .
echo ""

echo "Creating listing for user $USER2_ID"
curl -s ${BASE_URL}/public-api/listings \
  -X POST \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":$USER2_ID,\"listing_type\":\"sale\",\"price\":350000}" | jq .
echo ""

echo "Creating another listing for user $USER1_ID"
curl -s ${BASE_URL}/public-api/listings \
  -X POST \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":$USER1_ID,\"listing_type\":\"rent\",\"price\":5200}" | jq .
echo ""

# Test 4: Get enriched listings
echo "4. Get All Listings - ENRICHED with User Data (GET /public-api/listings)"
curl -s "${BASE_URL}/public-api/listings?page_num=1&page_size=10" | jq .
echo ""

# Test 5: Filter listings by user
echo "5. Get Listings by User $USER1_ID (GET /public-api/listings?user_id=$USER1_ID)"
curl -s "${BASE_URL}/public-api/listings?user_id=$USER1_ID" | jq .
echo ""

# Test 6: Test pagination
echo "6. Get Listings - Page 1 with Page Size 2"
curl -s "${BASE_URL}/public-api/listings?page_num=1&page_size=2" | jq .
echo ""

# Test 7: Error handling - Invalid JSON
echo "7. Error Handling - Invalid JSON Request Body"
curl -s ${BASE_URL}/public-api/users \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{invalid json}' | jq .
echo ""

# Test 8: Error handling - Missing required field
echo "8. Error Handling - Empty Name"
curl -s ${BASE_URL}/public-api/users \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":""}' | jq .
echo ""

echo "=================================="
echo "API Tests Completed"
echo "=================================="
echo ""
echo "Notice how GET /public-api/listings returns"
echo "listings with EMBEDDED user information!"
echo "This demonstrates the API Gateway pattern."
