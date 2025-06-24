#!/bin/bash

# Real-world scenario test script for schema-server with Keploy matcher logic

BASE_URL="http://localhost:7080"

# Real-world E-commerce API Schema (Base schema)
ECOMMERCE_SCHEMA='{
  "openapi": "3.0.0",
  "info": {
    "title": "E-commerce API",
    "version": "1.0.0",
    "description": "Production e-commerce API for online store"
  },
  "paths": {
    "/products": {
      "get": {
        "summary": "Get all products",
        "operationId": "getProducts",
        "parameters": [
          {
            "name": "category",
            "in": "query",
            "description": "Filter by product category",
            "required": false,
            "schema": { "type": "string" },
            "example": "electronics"
          },
          {
            "name": "limit",
            "in": "query",
            "description": "Number of products to return",
            "required": false,
            "schema": { "type": "integer", "minimum": 1, "maximum": 100 },
            "example": 20
          }
        ],
        "responses": {
          "200": {
            "description": "List of products",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "products": {
                      "type": "array",
                      "items": { "$ref": "#/components/schemas/Product" }
                    },
                    "total": { "type": "integer" },
                    "page": { "type": "integer" }
                  }
                },
                "example": {
                  "products": [
                    {
                      "id": "prod-123",
                      "name": "Wireless Headphones",
                      "price": 199.99,
                      "category": "electronics",
                      "in_stock": true
                    }
                  ],
                  "total": 150,
                  "page": 1
                }
              }
            }
          },
          "400": {
            "description": "Bad request",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Error" },
                "example": { "code": "INVALID_CATEGORY", "message": "Invalid category specified" }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create a new product",
        "operationId": "createProduct",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/ProductCreate" },
              "example": {
                "name": "Smart Watch",
                "price": 299.99,
                "category": "electronics",
                "description": "Advanced fitness tracking smartwatch"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Product created successfully",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Product" },
                "example": {
                  "id": "prod-124",
                  "name": "Smart Watch",
                  "price": 299.99,
                  "category": "electronics",
                  "in_stock": true
                }
              }
            }
          }
        }
      }
    },
    "/orders": {
      "post": {
        "summary": "Create a new order",
        "operationId": "createOrder",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/OrderCreate" },
              "example": {
                "customer_id": "cust-456",
                "items": [
                  { "product_id": "prod-123", "quantity": 2 },
                  { "product_id": "prod-124", "quantity": 1 }
                ],
                "shipping_address": {
                  "street": "123 Main St",
                  "city": "San Francisco",
                  "state": "CA",
                  "zip": "94105"
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Order created successfully",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Order" },
                "example": {
                  "id": "ord-789",
                  "status": "pending",
                  "total": 699.97,
                  "created_at": "2023-12-01T10:30:00Z"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Product": {
        "type": "object",
        "required": ["id", "name", "price", "category"],
        "properties": {
          "id": { "type": "string", "description": "Unique product identifier" },
          "name": { "type": "string", "description": "Product name" },
          "price": { "type": "number", "format": "decimal", "minimum": 0 },
          "category": { "type": "string", "enum": ["electronics", "clothing", "home", "books"] },
          "description": { "type": "string" },
          "in_stock": { "type": "boolean", "default": true }
        }
      },
      "ProductCreate": {
        "type": "object",
        "required": ["name", "price", "category"],
        "properties": {
          "name": { "type": "string" },
          "price": { "type": "number", "format": "decimal", "minimum": 0 },
          "category": { "type": "string", "enum": ["electronics", "clothing", "home", "books"] },
          "description": { "type": "string" }
        }
      },
      "OrderCreate": {
        "type": "object",
        "required": ["customer_id", "items"],
        "properties": {
          "customer_id": { "type": "string" },
          "items": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "product_id": { "type": "string" },
                "quantity": { "type": "integer", "minimum": 1 }
              }
            }
          },
          "shipping_address": { "$ref": "#/components/schemas/Address" }
        }
      },
      "Order": {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "status": { "type": "string", "enum": ["pending", "processing", "shipped", "delivered"] },
          "total": { "type": "number", "format": "decimal" },
          "created_at": { "type": "string", "format": "date-time" }
        }
      },
      "Address": {
        "type": "object",
        "required": ["street", "city", "state", "zip"],
        "properties": {
          "street": { "type": "string" },
          "city": { "type": "string" },
          "state": { "type": "string" },
          "zip": { "type": "string" }
        }
      },
      "Error": {
        "type": "object",
        "required": ["code", "message"],
        "properties": {
          "code": { "type": "string" },
          "message": { "type": "string" }
        }
      }
    }
  }
}'

# Identical schema (should match 100%)
IDENTICAL_SCHEMA="$ECOMMERCE_SCHEMA"



# Different schema (should not match)
DIFFERENT_SCHEMA='{
  "openapi": "3.0.0",
  "info": {
    "title": "Banking API",
    "version": "2.0.0",
    "description": "Core banking system API"
  },
  "paths": {
    "/accounts": {
      "get": {
        "summary": "Get user accounts",
        "operationId": "getAccounts",
        "parameters": [
          {
            "name": "user_id",
            "in": "query",
            "required": true,
            "schema": { "type": "string" },
            "example": "user-123"
          }
        ],
        "responses": {
          "200": {
            "description": "List of accounts",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "accounts": {
                      "type": "array",
                      "items": { "$ref": "#/components/schemas/Account" }
                    }
                  }
                },
                "example": {
                  "accounts": [
                    {
                      "id": "acc-456",
                      "type": "checking",
                      "balance": 1250.75,
                      "currency": "USD"
                    }
                  ]
                }
              }
            }
          }
        }
      }
    },
    "/transactions": {
      "post": {
        "summary": "Create a transaction",
        "operationId": "createTransaction",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/TransactionCreate" },
              "example": {
                "from_account": "acc-456",
                "to_account": "acc-789",
                "amount": 100.00,
                "description": "Payment for services"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Transaction created",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Transaction" },
                "example": {
                  "id": "txn-123",
                  "status": "completed",
                  "timestamp": "2023-12-01T15:30:00Z"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Account": {
        "type": "object",
        "required": ["id", "type", "balance"],
        "properties": {
          "id": { "type": "string" },
          "type": { "type": "string", "enum": ["checking", "savings", "credit"] },
          "balance": { "type": "number", "format": "decimal" },
          "currency": { "type": "string", "default": "USD" }
        }
      },
      "TransactionCreate": {
        "type": "object",
        "required": ["from_account", "to_account", "amount"],
        "properties": {
          "from_account": { "type": "string" },
          "to_account": { "type": "string" },
          "amount": { "type": "number", "format": "decimal", "minimum": 0.01 },
          "description": { "type": "string" }
        }
      },
      "Transaction": {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "status": { "type": "string", "enum": ["pending", "completed", "failed"] },
          "timestamp": { "type": "string", "format": "date-time" }
        }
      }
    }
  }
}'

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

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

print_comparison_result() {
    local test_name="$1"
    local response="$2"
    local expected_match="$3"
    
    echo -e "\n${PURPLE}Keploy Matcher Test: $test_name${NC}"
    echo "Response: $response"
    
    # Extract match and score from response
    if [[ $response == *'"match":'* ]]; then
        match_value=$(echo "$response" | grep -o '"match":[^,}]*' | cut -d':' -f2 | tr -d ' ')
        score_value=$(echo "$response" | grep -o '"score":[^,}]*' | cut -d':' -f2 | tr -d ' ')
        
        echo "Match: $match_value, Score: $score_value"
        
        if [[ $match_value == "$expected_match" ]]; then
            echo -e "${GREEN}✅ PASSED - Match result as expected${NC}"
        else
            echo -e "${RED}❌ FAILED - Expected match: $expected_match, Got: $match_value${NC}"
        fi
    else
        echo -e "${RED}❌ FAILED - Invalid response format${NC}"
    fi
    echo "----------------------------------------"
}

echo -e "${YELLOW}🚀 Real-World Schema Server Test with Keploy Matcher Logic${NC}"

# 1. Health Check
echo -e "\n${YELLOW}1. Health Check${NC}"
response=$(curl -s $BASE_URL/api/v1/health)
print_result "Health Check" "$response" "healthy"

# 2. Upload Base E-commerce Schema
echo -e "\n${YELLOW}2. Upload Base E-commerce Schema${NC}"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas \
  -H "Content-Type: application/json" \
  -d "$ECOMMERCE_SCHEMA")
print_result "Upload E-commerce Schema" "$response" "uploaded successfully"

# 3. List Schemas
echo -e "\n${YELLOW}3. List Schemas${NC}"
response=$(curl -s $BASE_URL/api/v1/schemas)
print_result "List Schemas" "$response" "E-commerce API"

# 4. Fetch Schema
echo -e "\n${YELLOW}4. Fetch Schema${NC}"
response=$(curl -s "$BASE_URL/api/v1/schemas/E-commerce%20API")
echo "$response" > fetched_schema.json
print_result "Fetch Schema" "$response" "E-commerce API"

# 5. Download Schema
echo -e "\n${YELLOW}5. Download Schema${NC}"
curl -s "$BASE_URL/api/v1/schemas/E-commerce%20API/download" -o downloaded_schema.json
if [ -f downloaded_schema.json ]; then
    echo -e "${GREEN}✅ Downloaded schema file exists${NC}"
else
    echo -e "${RED}❌ Downloaded schema file not found${NC}"
fi

# 6. Keploy Matcher Tests - Real World Scenarios
echo -e "\n${YELLOW}6. Keploy Matcher Logic Tests${NC}"

# Test 1: Identical schema comparison (should match)
echo -e "\n${PURPLE}Test 1: Identical Schema Comparison${NC}"
compare_response=$(curl -s -X POST "$BASE_URL/api/v1/schemas/E-commerce%20API/compare" \
  -H "Content-Type: application/json" \
  -d "$IDENTICAL_SCHEMA")
print_comparison_result "Identical Schema vs Stored Schema" "$compare_response" "true"



# Test 2: Different schema comparison (should not match)
echo -e "\n${PURPLE}Test 2: Different Schema Comparison (Banking vs E-commerce)${NC}"
compare_response=$(curl -s -X POST "$BASE_URL/api/v1/schemas/E-commerce%20API/compare" \
  -H "Content-Type: application/json" \
  -d "$DIFFERENT_SCHEMA")
print_comparison_result "Banking Schema vs E-commerce Schema" "$compare_response" "false"

# Test 3: Upload Banking schema and test cross-comparison
echo -e "\n${PURPLE}Test 3: Upload Banking Schema and Cross-Compare${NC}"
response=$(curl -s -X POST $BASE_URL/api/v1/schemas \
  -H "Content-Type: application/json" \
  -d "$DIFFERENT_SCHEMA")
print_result "Upload Banking Schema" "$response" "uploaded successfully"

# Compare E-commerce vs Banking
compare_response=$(curl -s -X POST "$BASE_URL/api/v1/schemas/Banking%20API/compare" \
  -H "Content-Type: application/json" \
  -d "$ECOMMERCE_SCHEMA")
print_comparison_result "E-commerce Schema vs Banking Schema" "$compare_response" "false"

# Test 4: Edge case - empty paths comparison
echo -e "\n${PURPLE}Test 4: Edge Case - Schema with Minimal Paths${NC}"
MINIMAL_SCHEMA='{
  "openapi": "3.0.0",
  "info": {
    "title": "E-commerce API",
    "version": "1.0.0",
    "description": "Minimal version"
  },
  "paths": {},
  "components": {}
}'

compare_response=$(curl -s -X POST "$BASE_URL/api/v1/schemas/E-commerce%20API/compare" \
  -H "Content-Type: application/json" \
  -d "$MINIMAL_SCHEMA")
print_comparison_result "Minimal Schema vs Full E-commerce Schema" "$compare_response" "false"

# Cleanup
echo -e "\n${YELLOW}Cleanup${NC}"
rm -f fetched_schema.json downloaded_schema.json

echo -e "\n${GREEN}🎉 Real-World Keploy Matcher Schema Server Test Complete!${NC}"
echo -e "${BLUE}📊 Summary:${NC}"
echo -e "• Tested identical schema matching (expect: perfect match)"
echo -e "• Tested different domain schemas (expect: no match)"
echo -e "• Tested cross-domain comparison (expect: no match)"
echo -e "• Tested edge cases with minimal schemas (expect: no match)"
echo -e "\n${PURPLE}🔍 Keploy matcher logic provides semantic OpenAPI comparison!${NC}" 