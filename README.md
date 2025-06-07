# Schema Server

A completely independent HTTP server built with Go's standard library that provides RESTful APIs for managing schemas. This server is self-contained with no external dependencies on the Keploy project, making it easy to deploy and use standalone.

## Features

- **Schema Generation**: Generate schemas for mocks and tests with mock implementations
- **Schema Download**: Download schemas with various configurations
- **Schema Validation**: Validate schemas against contracts
- **HTTPDoc to OpenAPI Conversion**: Convert HTTPDoc format to OpenAPI specification
- **Health Check**: Monitor server health status
- **Mock Data**: Returns realistic mock data for all operations

## Architecture

The server is completely self-contained with internal implementations of all required types and interfaces:

### Directory Structure

```
schema-server/
├── main.go                    # Application entry point
├── server/                    # HTTP server implementation
│   └── server.go             # Server routes and handlers
├── internal/                  # Internal packages (no external dependencies)
│   ├── config/               # Configuration types
│   │   └── config.go
│   ├── models/               # Data models and types
│   │   └── models.go
│   └── contract/             # Contract service interface and mock implementation
│       └── service.go
├── go.mod                     # Go module definition (minimal dependencies)
├── test_server.sh            # Test script for all endpoints
└── README.md                 # This documentation
```

### Components

#### Internal Packages

- **`internal/config`**: Configuration structures for server settings
- **`internal/models`**: All data models including HTTPDoc, OpenAPI, TestCase, etc.
- **`internal/contract`**: Service interface and mock implementation with realistic responses

#### Server Package (`server/server.go`)

- **Server struct**: Main server structure with HTTP multiplexer, logger, and services
- **Route Handlers**: Individual handlers for each API endpoint with full functionality
- **Helper Functions**: JSON response utilities and comprehensive error handling

#### Main Package (`main.go`)

- **Initialization**: Logger setup and configuration loading
- **Service Creation**: Mock contract service initialization
- **Server Startup**: HTTP server creation and startup

## API Endpoints

### Health Check

- **GET** `/api/v1/health` - Check server health status

### Schema Operations

#### Fetching Schemas

- **GET** `/api/v1/schemas/tests` - Get all tests schema (returns mock data)
- **GET** `/api/v1/schemas/mocks/downloaded` - Get all downloaded mocks schemas (returns mock data)

#### Generating Schemas

- **POST** `/api/v1/schemas/generate` - Generate schemas (general)

  ```json
  {
    "checkConfig": true
  }
  ```

- **POST** `/api/v1/schemas/mocks/generate` - Generate mocks schemas

  ```json
  {
    "services": ["service1", "service2"],
    "mappings": {
      "service1": ["url1", "url2"],
      "service2": ["url3", "url4"]
    }
  }
  ```

- **POST** `/api/v1/schemas/tests/generate` - Generate tests schemas
  ```json
  {
    "selectedTests": ["test1", "test2"]
  }
  ```

#### Downloading Schemas

- **POST** `/api/v1/schemas/download` - Download schemas

  ```json
  {
    "checkConfig": true
  }
  ```

- **POST** `/api/v1/schemas/tests/download` - Download tests

  ```json
  {
    "path": "/path/to/tests"
  }
  ```

- **POST** `/api/v1/schemas/mocks/download` - Download mocks
  ```json
  {
    "path": "/path/to/mocks"
  }
  ```

#### Validation

- **POST** `/api/v1/schemas/validate` - Validate schemas

#### Conversion

- **POST** `/api/v1/schemas/convert` - Convert HTTPDoc to OpenAPI (returns actual converted OpenAPI spec)
  ```json
  {
    "name": "Example API",
    "version": "1.0.0",
    "kind": "HTTPDoc",
    "spec": {
      "request": {
        "method": "GET",
        "url": "https://api.example.com/users",
        "header": {
          "Authorization": "Bearer token"
        },
        "body": ""
      },
      "response": {
        "statusCode": 200,
        "statusMessage": "OK",
        "body": "{\"users\": []}"
      }
    }
  }
  ```

## Setup and Running

### Prerequisites

- Go 1.22.0 or later
- **No external dependencies** - completely self-contained!

### Installation

1. Navigate to the schema-server directory:

   ```bash
   cd schema-server
   ```

2. Download dependencies (only zap for logging):

   ```bash
   go mod tidy
   ```

3. Set environment variables (optional):

   ```bash
   export PORT=8080  # Default port if not set
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

The server will start on port 8080 (or the port specified in the PORT environment variable).

### Testing

Run the comprehensive test script:

```bash
chmod +x test_server.sh
./test_server.sh
```

Or test individual endpoints:

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Get tests schema (returns mock data)
curl http://localhost:8080/api/v1/schemas/tests

# Generate schemas
curl -X POST http://localhost:8080/api/v1/schemas/generate \
  -H "Content-Type: application/json" \
  -d '{"checkConfig": true}'

# Convert HTTPDoc to OpenAPI
curl -X POST http://localhost:8080/api/v1/schemas/convert \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test API",
    "version": "1.0.0",
    "kind": "HTTPDoc",
    "spec": {
      "request": {
        "method": "GET",
        "url": "https://api.test.com/users",
        "header": {},
        "body": ""
      },
      "response": {
        "statusCode": 200,
        "statusMessage": "OK",
        "body": "{\"users\": []}"
      }
    }
  }'
```

## Dependencies

The server has **minimal dependencies**:

- `go.uber.org/zap` - Structured logging (only external dependency)
- Go standard library - HTTP server, JSON handling, etc.

**No Keploy dependencies** - completely independent!

## Mock Implementation

The server includes a fully functional mock implementation that:

- Logs all operations with structured logging
- Returns realistic mock data for schema operations
- Simulates actual OpenAPI schema generation and conversion
- Provides proper error handling and validation
- Supports all API endpoints with meaningful responses

## Key Features

✅ **Zero External Dependencies**: No dependency on Keploy codebase  
✅ **Complete Implementation**: All endpoints return proper responses  
✅ **Realistic Mock Data**: Returns actual OpenAPI schemas, not placeholders  
✅ **Full Error Handling**: Comprehensive error responses and validation  
✅ **Structured Logging**: Detailed logging for all operations  
✅ **Standard HTTP**: Uses Go's standard library for maximum compatibility  
✅ **Easy Testing**: Includes test script for all endpoints  
✅ **Self-Contained**: Can be deployed anywhere without external services

## Example Responses

### Health Check Response

```json
{
  "status": "healthy",
  "service": "schema-server"
}
```

### Get Tests Schema Response

```json
{
  "data": {
    "test-set-1": {
      "test-case-1": {
        "openapi": "3.0.0",
        "info": {
          "title": "Mock Test API",
          "version": "1.0.0"
        },
        "paths": {
          "/test": {
            "get": {
              "summary": "Mock test endpoint",
              "responses": {
                "200": {
                  "description": "Success"
                }
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

### HTTPDoc to OpenAPI Conversion Response

```json
{
  "data": {
    "openapi": "3.0.0",
    "info": {
      "title": "Test API",
      "version": "1.0.0",
      "description": "HTTPDoc"
    },
    "servers": [
      {
        "url": "https://api.example.com"
      }
    ],
    "paths": {
      "/converted": {
        "get": {
          "summary": "Converted from HTTPDoc",
          "description": "Generated from Test API",
          "responses": {
            "200": {
              "description": "Successful response",
              "content": {
                "application/json": {
                  "schema": {
                    "type": "object"
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

## Future Enhancements

1. **Real Database Integration**: Replace mock service with actual database backends
2. **Authentication**: Add JWT or API key authentication
3. **Rate Limiting**: Implement request rate limiting
4. **CORS Support**: Add CORS headers for web client support
5. **Configuration Files**: Support for YAML/JSON configuration files
6. **Metrics**: Add Prometheus metrics for monitoring
7. **Docker Support**: Add Dockerfile for containerized deployment
8. **File Storage**: Add file-based schema storage and retrieval
