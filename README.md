# Schema Server

A production-ready HTTP server built with Go that provides comprehensive RESTful APIs for **advanced schema generation and management**. This server features sophisticated schema generation capabilities with file-based persistence, making it ideal for API contract testing and documentation workflows.

## 🚀 Key Features

- **🔄 Advanced Schema Generation**: Sophisticated real-time schema generation from HTTP traffic and test data
- **📁 Organized File Storage**: Structured schema storage in organized directory hierarchies
- **🔍 Intelligent Type Inference**: Automatic OpenAPI schema generation with smart JSON type detection
- **📥 Flexible Export**: Download and export schemas to custom locations with multiple formats
- **✅ Comprehensive Validation**: Multi-layered schema validation with detailed error reporting
- **🔄 Format Conversion**: Convert between HTTPDoc and OpenAPI specifications seamlessly
- **🏥 Production Monitoring**: Built-in health checks and operational endpoints
- **🧪 Extensive Testing**: Comprehensive test suites covering all functionality

## 🏗️ Architecture

The server implements a modular architecture with sophisticated schema generation capabilities:

### Directory Structure

```
schema-server/
├── main.go                    # Application entry point
├── http/                      # HTTP server implementation
│   └── server.go             # Server routes and handlers
├── internal/                  # Internal packages
│   ├── models/               # Data models and types
│   │   └── models.go
│   └── contract/             # Schema generation core
│       ├── interfaces.go     # Database and storage interfaces
│       ├── generator.go      # Main schema generator
│       ├── openapi_converter.go  # OpenAPI conversion logic
│       ├── schema_operations.go  # Schema operations (mocks/tests)
│       ├── validation.go     # Schema validation
│       ├── schema_manager.go # Service manager
│       └── file_operations.go    # File I/O operations
├── schemas/                   # Generated schema files (auto-created)
│   ├── tests/                # Test schemas by service
│   ├── mocks/                # Mock schemas by service
│   └── converted/            # Converted schema files
├── go.mod                     # Go module dependencies
├── test_server.sh            # Comprehensive test script
└── README.md                 # Documentation
```

### 🔧 Core Components

#### Schema Generator (`internal/contract/`)

- **OpenAPI Converter**: Advanced HTTPDoc to OpenAPI conversion with intelligent type inference
- **Schema Operations**: Comprehensive mock and test schema generation with service mapping
- **Validation Engine**: Multi-level schema validation with detailed error reporting
- **File Operations**: Organized file I/O with automatic directory management
- **Interface Abstractions**: Clean database and storage interfaces for extensibility

## 🌐 API Endpoints

All endpoints run on **port 7080** by default.

### Health Check

- **GET** `/api/v1/health` - Check server health status

### Schema Operations

#### Fetching Schemas

- **GET** `/api/v1/schemas/tests` - Get all test schemas with organized structure
- **GET** `/api/v1/schemas/mocks/downloaded` - Get all downloaded mock schemas by service

#### Generating Schemas

- **POST** `/api/v1/schemas/generate` - Generate comprehensive schemas

  ```json
  {
    "checkConfig": true
  }
  ```

- **POST** `/api/v1/schemas/mocks/generate` - Generate mock schemas from service endpoints

  ```json
  {
    "services": ["user-service", "order-service"],
    "mappings": {
      "user-service": ["test-set-1", "test-set-2"],
      "order-service": ["test-set-3", "test-set-4"]
    }
  }
  ```

- **POST** `/api/v1/schemas/tests/generate` - Generate test schemas from test cases
  ```json
  {
    "selectedTests": ["test-set-1", "test-set-2", "test-set-3"]
  }
  ```

#### Downloading Schemas

- **POST** `/api/v1/schemas/download` - Download schemas from external sources

  ```json
  {
    "checkConfig": true
  }
  ```

- **POST** `/api/v1/schemas/tests/download` - Export test schemas to specified location

  ```json
  {
    "path": "/tmp/exported-tests"
  }
  ```

- **POST** `/api/v1/schemas/mocks/download` - Export mock schemas to specified location
  ```json
  {
    "path": "/tmp/exported-mocks"
  }
  ```

#### Validation & Conversion

- **POST** `/api/v1/schemas/validate` - Validate all stored schemas

- **POST** `/api/v1/schemas/{title}/compare` - Compare schemas using Keploy matcher logic

  ```json
  {
    "openapi": "3.0.0",
    "info": {
      "title": "API to compare",
      "version": "1.0.0"
    },
    "paths": {
      "/example": {
        "get": {
          "responses": {
            "200": {
              "description": "Success"
            }
          }
        }
      }
    }
  }
  ```

- **POST** `/api/v1/schemas/convert` - Convert HTTPDoc to OpenAPI with intelligent type inference
  ```json
  {
    "name": "User Management API",
    "version": "1.0.0",
    "kind": "HTTPDoc",
    "spec": {
      "request": {
        "method": "POST",
        "url": "https://api.example.com/users",
        "header": {
          "Content-Type": "application/json",
          "Authorization": "Bearer token"
        },
        "body": "{\"name\": \"John Doe\", \"email\": \"john@example.com\", \"age\": 30}"
      },
      "response": {
        "statusCode": 201,
        "statusMessage": "Created",
        "body": "{\"id\": 123, \"name\": \"John Doe\", \"created_at\": \"2024-01-15T10:30:00Z\"}"
      }
    }
  }
  ```

## 🚀 Setup and Running

### Prerequisites

- Go 1.22.0 or later
- Dependencies managed via `go.mod` (automatically installed)

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
   export SCHEMAS_PATH=./schemas       # Default schemas directory
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

The server will start on **port 7080** and automatically create the `./schemas/` directory structure.

## 🧪 Testing

### **1. Comprehensive Test Suite (Recommended)**

The test script covers all functionality with realistic data scenarios:

```bash
# Make executable and run
chmod +x test.sh
./test.sh
```

**Test Coverage:**

- Health check validation
- Schema generation and file persistence
- Mock schema generation with service mappings
- Test schema generation from test cases
- Schema export to custom locations
- HTTPDoc to OpenAPI conversion with intelligent type inference
- Identical schema matching verification
- Different domain schema comparison (banking vs e-commerce)
- Cross-domain schema validation
- Edge case testing with minimal schemas
- Error handling and validation
- File system verification

### **2. Manual Testing**

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
- ✅ Keploy matcher logic for schema comparison
- ✅ Identical schema matching (perfect score)
- ✅ Cross-domain schema comparison (banking vs e-commerce)
- ✅ Edge case handling with minimal schemas
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

You can also test individual endpoints manually:

### **Manual Testing Examples**

```bash
# Health check
curl http://localhost:7080/api/v1/health

# Generate schemas with real data
curl -X POST http://localhost:7080/api/v1/schemas/generate \
  -H "Content-Type: application/json" \
  -d '{"checkConfig": true}'

# Generate service mocks from test data
curl -X POST http://localhost:7080/api/v1/schemas/mocks/generate \
  -H "Content-Type: application/json" \
  -d '{
    "services": ["payment-service", "user-service"],
    "mappings": {
      "payment-service": ["test-set-1", "test-set-2"],
      "user-service": ["test-set-3", "test-set-4"]
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
schemas/
├── converted/
│   └── User Management API-converted.yaml    # HTTPDoc conversions
├── mocks/
│   ├── payment-service/
│   │   ├── test-set-1/
│   │   │   ├── mock-1.yaml
│   │   │   └── mock-2.yaml
│   │   └── test-set-2/
│   │       └── mock-1.yaml
│   ├── user-service/
│   │   └── test-set-3/
│   │       └── mock-1.yaml
│   └── notification-service/
│       └── test-set-4/
│           └── mock-1.yaml
└── tests/
    ├── test-set-1/
    │   ├── test-case-1.yaml                   # Generated test schemas
    │   └── test-case-2.yaml
    ├── test-set-2/
    │   └── test-case-1.yaml
    └── test-set-3/
        └── test-case-1.yaml
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

**Completely independent** - no external framework dependencies!

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
