#!/bin/bash

echo "🚀 Testing Schema Server with Real Data"
echo "======================================="

BASE_URL="http://localhost:7080"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print test results
print_result() {
    local test_name="$1"
    local response="$2"
    local expected="$3"
    
    echo -e "\n${BLUE}Testing: $test_name${NC}"
    echo "Response: $response"
    
    if [[ $response == *"$expected"* ]]; then
        echo -e "${GREEN}✅ PASSED${NC}"
    else
        echo -e "${RED}❌ FAILED${NC}"
    fi
    echo "----------------------------------------"
}

# 1. Health Check
echo -e "\n${YELLOW}1. Health Check${NC}"
response=$(curl -s $BASE_URL/api/v1/health)
print_result "Health Check" "$response" "healthy"

# 2. Generate Initial Schemas
echo -e "\n${YELLOW}2. Generate Initial Schemas${NC}"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas/generate \
  -H "Content-Type: application/json" \
  -d '{"checkConfig": true}')
print_result "Generate Schemas" "$response" "successfully"

# 3. Get All Tests Schema (should now return real data from files)
echo -e "\n${YELLOW}3. Get All Tests Schema${NC}"
response=$(curl -s $BASE_URL/api/v1/schemas/tests)
print_result "Get Tests Schema" "$response" "User Service Test"

# 4. Get Downloaded Mocks Schema
echo -e "\n${YELLOW}4. Get Downloaded Mocks Schema${NC}"
response=$(curl -s $BASE_URL/api/v1/schemas/mocks/downloaded)
print_result "Get Mocks Schema" "$response" "Payment Service Mock"

# 5. Generate Mock Schemas with Real Data
echo -e "\n${YELLOW}5. Generate Mock Schemas with Real Service Data${NC}"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas/mocks/generate \
  -H "Content-Type: application/json" \
  -d '{
    "services": ["user-service", "inventory-service", "order-service"],
    "mappings": {
      "user-service": [
        "https://api.example.com/users",
        "https://api.example.com/users/{id}",
        "https://api.example.com/users/{id}/profile"
      ],
      "inventory-service": [
        "/inventory/items",
        "/inventory/items/{id}",
        "/inventory/stock"
      ],
      "order-service": [
        "/orders",
        "/orders/{id}",
        "/orders/{id}/status"
      ]
    }
  }')
print_result "Generate Mock Schemas" "$response" "successfully"

# 6. Generate Test Schemas with Real Test Names
echo -e "\n${YELLOW}6. Generate Test Schemas with Real Test Names${NC}"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas/tests/generate \
  -H "Content-Type: application/json" \
  -d '{
    "selectedTests": [
      "user-authentication-test",
      "payment-processing-test",
      "order-fulfillment-test",
      "inventory-update-test"
    ]
  }')
print_result "Generate Test Schemas" "$response" "successfully"

# 7. Download Schemas
echo -e "\n${YELLOW}7. Download Schemas${NC}"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas/download \
  -H "Content-Type: application/json" \
  -d '{"checkConfig": true}')
print_result "Download Schemas" "$response" "successfully"

# 8. Download Tests to Specific Path
echo -e "\n${YELLOW}8. Download Tests to Specific Path${NC}"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas/tests/download \
  -H "Content-Type: application/json" \
  -d '{"path": "/tmp/downloaded-tests"}')
print_result "Download Tests" "$response" "successfully"

# 9. Download Mocks to Specific Path
echo -e "\n${YELLOW}9. Download Mocks to Specific Path${NC}"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas/mocks/download \
  -H "Content-Type: application/json" \
  -d '{"path": "/tmp/downloaded-mocks"}')
print_result "Download Mocks" "$response" "successfully"

# 10. Validate All Schemas
echo -e "\n${YELLOW}10. Validate All Schemas${NC}"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas/validate)
print_result "Validate Schemas" "$response" "successfully"

# 11. Convert HTTPDoc to OpenAPI with Real Data
echo -e "\n${YELLOW}11. Convert HTTPDoc to OpenAPI with Real Data${NC}"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas/convert \
  -H "Content-Type: application/json" \
  -d '{
    "name": "E-commerce User API",
    "version": "2.1.0",
    "kind": "HTTPDoc",
    "spec": {
      "request": {
        "method": "POST",
        "url": "https://api.ecommerce.com/users",
        "header": {
          "Content-Type": "application/json",
          "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
        },
        "body": "{\"name\": \"John Doe\", \"email\": \"john@example.com\", \"age\": 30, \"preferences\": {\"newsletter\": true, \"theme\": \"dark\"}}"
      },
      "response": {
        "statusCode": 201,
        "statusMessage": "Created",
        "header": {
          "Content-Type": "application/json",
          "Location": "/users/12345"
        },
        "body": "{\"id\": 12345, \"name\": \"John Doe\", \"email\": \"john@example.com\", \"created_at\": \"2024-01-15T10:30:00Z\", \"status\": \"active\"}"
      }
    }
  }')
print_result "Convert HTTPDoc to OpenAPI" "$response" "E-commerce User API"

# 12. Generate Complex Mock Schemas
echo -e "\n${YELLOW}12. Generate Complex Mock Schemas${NC}"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas/mocks/generate \
  -H "Content-Type: application/json" \
  -d '{
    "services": ["payment-gateway", "notification-service"],
    "mappings": {
      "payment-gateway": [
        "https://gateway.payment.com/v2/charges",
        "https://gateway.payment.com/v2/refunds",
        "https://gateway.payment.com/v2/webhooks"
      ],
      "notification-service": [
        "/notifications/email",
        "/notifications/sms", 
        "/notifications/push"
      ]
    }
  }')
print_result "Generate Complex Mock Schemas" "$response" "successfully"

# 13. Test Error Handling
echo -e "\n${YELLOW}13. Test Error Handling${NC}"
echo "Testing invalid HTTP method:"
response=$(curl -s -X DELETE $BASE_URL/api/v1/health)
print_result "Invalid Method Error" "$response" "Method not allowed"

echo "Testing invalid JSON:"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas/generate \
  -H "Content-Type: application/json" \
  -d '{"invalid": json}')
print_result "Invalid JSON Error" "$response" "Invalid request body"

echo "Testing invalid endpoint:"
response=$(curl -s $BASE_URL/api/v1/invalid-endpoint)
print_result "Invalid Endpoint" "$response" "404"

# 14. Check if files were actually created
echo -e "\n${YELLOW}14. Verify File Creation${NC}"
if [ -d "./contracts" ]; then
    echo -e "${GREEN}✅ Contracts directory exists${NC}"
    echo "Directory structure:"
    find ./contracts -type f -name "*.json" | head -10
else
    echo -e "${RED}❌ Contracts directory not found${NC}"
fi

# 15. Final Status Check
echo -e "\n${YELLOW}15. Final Status Check${NC}"
response=$(curl -s $BASE_URL/api/v1/health)
print_result "Final Health Check" "$response" "healthy"

echo -e "\n${GREEN}🎉 Testing Complete!${NC}"
echo -e "${BLUE}Check the ./contracts directory to see generated schema files${NC}"
echo -e "${BLUE}Check /tmp/downloaded-tests and /tmp/downloaded-mocks for downloaded files${NC}" 