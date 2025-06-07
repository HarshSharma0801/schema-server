#!/bin/bash

echo "Testing Schema Server..."
echo "========================"

BASE_URL="http://localhost:7080"

# Test health check
echo -e "\n1. Testing Health Check:"
curl -s $BASE_URL/api/v1/health | jq '.'

# Test GET endpoints
echo -e "\n2. Testing GET endpoints:"
echo "Getting all tests schema:"
curl -s $BASE_URL/api/v1/schemas/tests | jq '.'

echo -e "\nGetting downloaded mocks:"
curl -s $BASE_URL/api/v1/schemas/mocks/downloaded | jq '.'

# Test POST endpoints
echo -e "\n3. Testing POST endpoints:"
echo "Generate schemas:"
curl -s -X POST $BASE_URL/api/v1/schemas/generate \
  -H "Content-Type: application/json" \
  -d '{"checkConfig": true}' | jq '.'

echo -e "\nValidate schemas:"
curl -s -X POST $BASE_URL/api/v1/schemas/validate | jq '.'

echo -e "\n4. Testing error handling:"
echo "Invalid method:"
curl -s -X DELETE $BASE_URL/api/v1/health | jq '.'

echo -e "\nInvalid endpoint:"
curl -s $BASE_URL/api/v1/invalid | jq '.'

echo -e "\nTesting complete!" 