# Schema Server

A completely independent HTTP server built with Go's standard library that provides RESTful APIs for **dynamic schema generation and management**. This server features **real file-based schema operations** with no hard-coded data, making it production-ready for actual schema management workflows.

## 🚀 Key Features

- **🔄 Dynamic Schema Generation**: Real-time schema generation from input data with file persistence
- **📁 File-Based Storage**: All schemas stored in organized directory structure (`./contracts/`)
- **🔍 Smart Schema Inference**: Automatic OpenAPI schema generation from JSON data and HTTPDoc
- **📥 Schema Download**: Download and export schemas to custom locations
- **✅ Schema Validation**: Validate all generated schemas for correctness
- **🔄 HTTPDoc to OpenAPI Conversion**: Convert HTTPDoc format to OpenAPI specification with schema inference
- **🏥 Health Monitoring**: Built-in health check endpoints
- **🧪 Comprehensive Testing**: Multiple test suites with real data scenarios

## 🏗️ Architecture

The server features **real schema generation** with file persistence and dynamic data processing:

### Directory Structure

```
schema-server/
├── main.go                    # Application entry point (Real Service)
├── server/                    # HTTP server implementation
│   └── server.go             # Server routes and handlers
├── internal/                  # Internal packages
│   ├── config/               # Configuration types
│   │   └── config.go
│   ├── models/               # Data models and types
│   │   └── models.go
│   └── contract/             # Contract service interface and REAL implementation
│       ├── service.go        # Service interface
│       └── real_service.go   # Real file-based implementation (NEW!)
├── contracts/                 # Generated schema files (created dynamically)
│   ├── tests/                # Test schemas organized by service
│   ├── mocks/                # Mock schemas organized by service
│   └── generated/            # Converted schemas and generated files
├── go.mod                     # Go module definition (minimal dependencies)
├── test_server.sh            # Original test script
├── simple_test.sh            # Simple test without jq dependency
├── real_test.sh              # Comprehensive test with real data scenarios (NEW!)
└── README.md                 # This documentation
```

### 🔧 Components

#### Real Service Implementation (`internal/contract/real_service.go`)

- **File-Based Operations**: Read/write schemas to disk with organized structure
- **Dynamic Generation**: Generate OpenAPI schemas from service endpoints and test names
- **Schema Inference**: Smart JSON-to-schema conversion with type detection
- **Directory Management**: Automatic creation and organization of schema directories
- **Download Operations**: Copy schemas to specified paths for distribution

## 🌐 API Endpoints

All endpoints run on **port 7080** by default.

### Health Check

- **GET** `/api/v1/health` - Check server health status

### Schema Operations

#### Fetching Schemas (Real Data from Files)

- **GET** `/api/v1/schemas/tests` - Get all tests schema (from `./contracts/tests/`)
- **GET** `/api/v1/schemas/mocks/downloaded` - Get all downloaded mocks schemas (from `./contracts/mocks/`)

#### Generating Schemas (Dynamic File Creation)

- **POST** `/api/v1/schemas/generate` - Generate schemas with file persistence

  ```json
  {
    "checkConfig": true
  }
  ```

- **POST** `/api/v1/schemas/mocks/generate` - Generate dynamic mock schemas from endpoints

  ```json
  {
    "services": ["user-service", "inventory-service"],
    "mappings": {
      "user-service": [
        "https://api.example.com/users",
        "https://api.example.com/users/{id}",
        "/users/{id}/profile"
      ],
      "inventory-service": ["/inventory/items", "/inventory/stock"]
    }
  }
  ```

- **POST** `/api/v1/schemas/tests/generate` - Generate test schemas dynamically
  ```json
  {
    "selectedTests": [
      "user-authentication-test",
      "payment-processing-test",
      "order-fulfillment-test"
    ]
  }
  ```

#### Downloading Schemas (File Operations)

- **POST** `/api/v1/schemas/download` - Download schemas from remote sources

  ```json
  {
    "checkConfig": true
  }
  ```

- **POST** `/api/v1/schemas/tests/download` - Download test files to custom location

  ```json
  {
    "path": "/tmp/downloaded-tests"
  }
  ```

- **POST** `/api/v1/schemas/mocks/download` - Download mock files to custom location
  ```json
  {
    "path": "/tmp/downloaded-mocks"
  }
  ```

#### Validation & Conversion

- **POST** `/api/v1/schemas/validate` - Validate all stored schemas

- **POST** `/api/v1/schemas/convert` - Convert HTTPDoc to OpenAPI with smart schema inference
  ```json
  {
    "name": "E-commerce User API",
    "version": "2.1.0",
    "kind": "HTTPDoc",
    "spec": {
      "request": {
        "method": "POST",
        "url": "https://api.ecommerce.com/users",
        "header": {
          "Content-Type": "application/json",
          "Authorization": "Bearer token123"
        },
        "body": "{\"name\": \"John Doe\", \"email\": \"john@example.com\", \"age\": 30}"
      },
      "response": {
        "statusCode": 201,
        "statusMessage": "Created",
        "body": "{\"id\": 12345, \"name\": \"John Doe\", \"created_at\": \"2024-01-15T10:30:00Z\"}"
      }
    }
  }
  ```

## 🚀 Setup and Running

### Prerequisites

- Go 1.22.0 or later
- **Only one external dependency**: `go.uber.org/zap` for logging

### Installation

1. Clone and navigate to the directory:

   ```bash
   cd schema-server
   ```

2. Download dependencies:

   ```bash
   go mod tidy
   ```

3. Set environment variables (optional):

   ```bash
   export PORT=7080                    # Default port
   export CONTRACTS_PATH=./contracts   # Default contracts directory
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

The server will start on **port 7080** and automatically create the `./contracts/` directory structure.

## 🧪 Testing

### **1. Comprehensive Test Suite (Recommended)**

Run the full test suite with real data scenarios:

```bash
chmod +x real_test.sh
./real_test.sh
```

This tests:

- ✅ Health checks
- ✅ Dynamic schema generation from real service data
- ✅ HTTPDoc to OpenAPI conversion with schema inference
- ✅ File creation and organization
- ✅ Download operations to custom paths
- ✅ Schema validation
- ✅ Error handling
- ✅ File verification

### **2. Simple Test Suite**

Test without `jq` dependency:

```bash
chmod +x simple_test.sh
./simple_test.sh
```

### **3. Original Test Suite**

Test with `jq` for formatted output (requires `brew install jq`):

```bash
chmod +x test_server.sh
./test_server.sh
```

### **4. Manual Testing Examples**

```bash
# Health check
curl http://localhost:7080/api/v1/health

# Generate schemas with real data
curl -X POST http://localhost:7080/api/v1/schemas/generate \
  -H "Content-Type: application/json" \
  -d '{"checkConfig": true}'

# Generate service mocks dynamically
curl -X POST http://localhost:7080/api/v1/schemas/mocks/generate \
  -H "Content-Type: application/json" \
  -d '{
    "services": ["payment-service", "user-service"],
    "mappings": {
      "payment-service": [
        "https://gateway.payment.com/v2/charges",
        "/refunds", "/webhooks"
      ],
      "user-service": [
        "/users", "/users/{id}", "/users/{id}/profile"
      ]
    }
  }'

# Convert HTTPDoc to OpenAPI with schema inference
curl -X POST http://localhost:7080/api/v1/schemas/convert \
  -H "Content-Type: application/json" \
  -d '{
    "name": "User Management API",
    "version": "1.0.0",
    "kind": "HTTPDoc",
    "spec": {
      "request": {
        "method": "POST",
        "url": "https://api.example.com/users",
        "body": "{\"name\": \"Alice\", \"email\": \"alice@example.com\", \"role\": \"admin\"}"
      },
      "response": {
        "statusCode": 201,
        "statusMessage": "Created",
        "body": "{\"id\": 42, \"name\": \"Alice\", \"created_at\": \"2024-01-15T10:30:00Z\"}"
      }
    }
  }'
```

## 📁 Generated File Structure

After running tests or operations, you'll see:

```
contracts/
├── generated/
│   └── E-commerce User API-converted.json    # HTTPDoc conversions
├── mocks/
│   ├── payment-service/
│   │   ├── mock-default.json
│   │   └── mock-20250609-121816.json         # Timestamped generations
│   ├── user-service/
│   │   └── mock-20250609-121816.json
│   └── notification-service/
│       └── mock-default.json
└── tests/
    ├── generated/
    │   ├── user-authentication-test.json      # Generated test schemas
    │   ├── payment-processing-test.json
    │   └── order-fulfillment-test.json
    ├── user-service/
    │   └── get-user-test.json                 # Service-specific tests
    └── order-service/
        └── create-order-test.json
```

## 📊 Example Responses

### Real Test Schema Response

```json
{
  "data": {
    "user-service": {
      "get-user-test": {
        "openapi": "3.0.0",
        "info": {
          "title": "User Service Test",
          "version": "1.0.0"
        },
        "paths": {
          "/users/{id}": {
            "get": {
              "summary": "Get user by ID",
              "responses": {
                "200": { "description": "User found" }
              }
            }
          }
        }
      }
    }
  },
  "message": "Tests schema retrieved successfully"
}
```

### Dynamic HTTPDoc Conversion

```json
{
  "data": {
    "openapi": "3.0.0",
    "info": {
      "title": "E-commerce User API",
      "version": "2.1.0",
      "description": "Generated from HTTPDoc"
    },
    "paths": {
      "/users": {
        "post": {
          "summary": "POST operation for E-commerce User API",
          "requestBody": {
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "name": { "type": "string" },
                    "email": { "type": "string" },
                    "age": { "type": "number" }
                  }
                }
              }
            }
          },
          "responses": {
            "201": {
              "description": "Created",
              "content": {
                "application/json": {
                  "schema": {
                    "type": "object",
                    "properties": {
                      "id": { "type": "number" },
                      "name": { "type": "string" },
                      "created_at": { "type": "string" }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "message": "HTTPDoc converted to OpenAPI successfully"
}
```

## 🔧 Key Improvements Over Mock Implementation

✅ **Real File Operations**: Schemas persisted to disk, not in-memory  
✅ **Dynamic Generation**: Creates schemas from actual input data  
✅ **Smart Schema Inference**: Analyzes JSON to generate OpenAPI schemas  
✅ **Organized Storage**: Hierarchical directory structure for easy management  
✅ **Timestamped Files**: Automatic versioning with timestamps  
✅ **Copy Operations**: Real file copying for downloads  
✅ **Validation**: Actual file-based schema validation  
✅ **Zero Hardcoding**: All data generated dynamically from inputs

## 🛠️ Dependencies

Minimal dependencies for maximum compatibility:

- `go.uber.org/zap` - Structured logging (only external dependency)
- Go standard library - HTTP server, JSON handling, file operations

**No Keploy dependencies** - completely independent!

## 🔮 Future Enhancements

1. **Database Integration**: Add PostgreSQL/MongoDB backends for large-scale storage
2. **Authentication**: JWT/API key authentication for production security
3. **Rate Limiting**: Request throttling for production deployment
4. **Advanced Schema Inference**: More sophisticated JSON-to-OpenAPI conversion
5. **Web UI**: Browser-based interface for schema management
6. **Docker Support**: Containerized deployment with Docker Compose
7. **Metrics & Monitoring**: Prometheus metrics and health dashboards
8. **Schema Versioning**: Git-based versioning for schema evolution
9. **Batch Operations**: Bulk schema generation and processing
10. **Real-time Updates**: WebSocket-based real-time schema synchronization
