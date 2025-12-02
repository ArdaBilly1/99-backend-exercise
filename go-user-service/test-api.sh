#!/bin/bash

# Test script for User Service API
# Make sure the service is running before executing this script

BASE_URL="http://localhost:7000"

echo "=================================="
echo "Testing User Service API"
echo "=================================="
echo ""

# Test 1: Health check
echo "1. Health Check (GET /users/ping)"
curl -s ${BASE_URL}/users/ping
echo -e "\n"

# Test 2: Create user
echo "2. Create User (POST /users)"
echo "Creating user: John Doe"
curl -s ${BASE_URL}/users -XPOST -d name="John Doe" | jq .
echo ""

echo "Creating user: Jane Smith"
curl -s ${BASE_URL}/users -XPOST -d name="Jane Smith" | jq .
echo ""

echo "Creating user: Bob Johnson"
curl -s ${BASE_URL}/users -XPOST -d name="Bob Johnson" | jq .
echo ""

# Test 3: Get specific user
echo "3. Get Specific User (GET /users/1)"
curl -s ${BASE_URL}/users/1 | jq .
echo ""

# Test 4: Get all users
echo "4. Get All Users (GET /users?page_num=1&page_size=10)"
curl -s "${BASE_URL}/users?page_num=1&page_size=10" | jq .
echo ""

# Test 5: Test pagination
echo "5. Get Users - Page 1 with Page Size 2"
curl -s "${BASE_URL}/users?page_num=1&page_size=2" | jq .
echo ""

# Test 6: Error handling - Invalid user ID
echo "6. Error Handling - Get Non-existent User (GET /users/999)"
curl -s ${BASE_URL}/users/999 | jq .
echo ""

# Test 7: Error handling - Empty name
echo "7. Error Handling - Create User with Empty Name"
curl -s ${BASE_URL}/users -XPOST -d name="" | jq .
echo ""

echo "=================================="
echo "API Tests Completed"
echo "=================================="
