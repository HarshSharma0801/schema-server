#!/bin/bash

echo "Testing Schema Server (Simple Version)..."
echo "========================================"

BASE_URL="http://localhost:7080"

# Test health check
echo -e "\n1. Health Check:"
response=$(curl -s $BASE_URL/api/v1/health)
if [[ $response == *"healthy"* ]]; then
    echo "✅ Health check passed"
else
    echo "❌ Health check failed: $response"
fi

# Test GET endpoints
echo -e "\n2. GET Endpoints:"
response=$(curl -s $BASE_URL/api/v1/schemas/tests)
if [[ $response == *"message"* ]]; then
    echo "✅ Get tests schema passed"
else
    echo "❌ Get tests schema failed"
fi

response=$(curl -s $BASE_URL/api/v1/schemas/mocks/downloaded)
if [[ $response == *"message"* ]]; then
    echo "✅ Get downloaded mocks passed"
else
    echo "❌ Get downloaded mocks failed"
fi

# Test POST endpoints
echo -e "\n3. POST Endpoints:"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas/generate \
  -H "Content-Type: application/json" \
  -d '{"checkConfig": true}')
if [[ $response == *"successfully"* ]]; then
    echo "✅ Generate schemas passed"
else
    echo "❌ Generate schemas failed"
fi

response=$(curl -s -X POST $BASE_URL/api/v1/schemas/validate)
if [[ $response == *"successfully"* ]]; then
    echo "✅ Validate schemas passed"
else
    echo "❌ Validate schemas failed"
fi

# Test error handling
echo -e "\n4. Error Handling:"
response=$(curl -s -X DELETE $BASE_URL/api/v1/health)
if [[ $response == *"Method not allowed"* ]]; then
    echo "✅ Error handling passed"
else
    echo "❌ Error handling failed"
fi

echo -e "\n✅ All basic tests completed!" 