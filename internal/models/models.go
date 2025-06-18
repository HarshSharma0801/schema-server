package models

import "time"

// HTTPDoc represents an HTTP document for schema operations
type HTTPDoc struct {
	Name    string  `json:"name" yaml:"name" bson:"name"`
	Version string  `json:"version" yaml:"version" bson:"version"`
	Kind    string  `json:"kind" yaml:"kind" bson:"kind"`
	Spec    APISpec `json:"spec" yaml:"spec" bson:"spec"`
}

// APISpec represents the specification of an API
type APISpec struct {
	Request  HTTPRequest  `json:"request" yaml:"request" bson:"request"`
	Response HTTPResponse `json:"response" yaml:"response" bson:"response"`
}

// HTTPRequest represents an HTTP request
type HTTPRequest struct {
	Method string            `json:"method" yaml:"method" bson:"method"`
	URL    string            `json:"url" yaml:"url" bson:"url"`
	Header map[string]string `json:"header" yaml:"header" bson:"header"`
	Body   string            `json:"body" yaml:"body" bson:"body"`
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	StatusCode    int               `json:"statusCode" yaml:"statusCode" bson:"status_code"`
	StatusMessage string            `json:"statusMessage" yaml:"statusMessage" bson:"status_message"`
	Header        map[string]string `json:"header" yaml:"header" bson:"header"`
	Body          string            `json:"body" yaml:"body" bson:"body"`
}

// OpenAPI represents an OpenAPI specification
type OpenAPI struct {
	OpenAPI    string                 `json:"openapi" yaml:"openapi" bson:"openapi"`
	Info       Info                   `json:"info" yaml:"info" bson:"info"`
	Servers    []map[string]string    `json:"servers,omitempty" yaml:"servers,omitempty" bson:"servers,omitempty"`
	Paths      map[string]PathItem    `json:"paths" yaml:"paths" bson:"paths"`
	Components map[string]interface{} `json:"components,omitempty" yaml:"components,omitempty" bson:"components,omitempty"`
}

// Info represents the info section of OpenAPI
type Info struct {
	Title       string `json:"title" yaml:"title" bson:"title"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" bson:"description,omitempty"`
	Version     string `json:"version" yaml:"version" bson:"version"`
}

// PathItem represents a path item in OpenAPI
type PathItem struct {
	Get    *Operation `json:"get,omitempty" yaml:"get,omitempty" bson:"get,omitempty"`
	Post   *Operation `json:"post,omitempty" yaml:"post,omitempty" bson:"post,omitempty"`
	Put    *Operation `json:"put,omitempty" yaml:"put,omitempty" bson:"put,omitempty"`
	Patch  *Operation `json:"patch,omitempty" yaml:"patch,omitempty" bson:"patch,omitempty"`
	Delete *Operation `json:"delete,omitempty" yaml:"delete,omitempty" bson:"delete,omitempty"`
}

// Operation represents an operation in OpenAPI
type Operation struct {
	Summary     string                  `json:"summary,omitempty" yaml:"summary,omitempty" bson:"summary,omitempty"`
	Description string                  `json:"description,omitempty" yaml:"description,omitempty" bson:"description,omitempty"`
	OperationID string                  `json:"operationId,omitempty" yaml:"operationId,omitempty" bson:"operation_id,omitempty"`
	Parameters  []Parameter             `json:"parameters,omitempty" yaml:"parameters,omitempty" bson:"parameters,omitempty"`
	RequestBody *RequestBody            `json:"requestBody,omitempty" yaml:"requestBody,omitempty" bson:"request_body,omitempty"`
	Responses   map[string]ResponseItem `json:"responses" yaml:"responses" bson:"responses"`
}

// Parameter represents a parameter in OpenAPI
type Parameter struct {
	Name     string      `json:"name" yaml:"name" bson:"name"`
	In       string      `json:"in" yaml:"in" bson:"in"`
	Required bool        `json:"required,omitempty" yaml:"required,omitempty" bson:"required,omitempty"`
	Schema   ParamSchema `json:"schema" yaml:"schema" bson:"schema"`
	Example  interface{} `json:"example,omitempty" yaml:"example,omitempty" bson:"example,omitempty"`
}

// ParamSchema represents a parameter schema
type ParamSchema struct {
	Type string `json:"type" yaml:"type" bson:"type"`
}

// RequestBody represents a request body in OpenAPI
type RequestBody struct {
	Content map[string]MediaType `json:"content" yaml:"content" bson:"content"`
}

// ResponseItem represents a response item in OpenAPI
type ResponseItem struct {
	Description string               `json:"description" yaml:"description" bson:"description"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty" bson:"content,omitempty"`
}

// MediaType represents a media type in OpenAPI
type MediaType struct {
	Schema  Schema      `json:"schema" yaml:"schema" bson:"schema"`
	Example interface{} `json:"example,omitempty" yaml:"example,omitempty" bson:"example,omitempty"`
}

// Schema represents a schema in OpenAPI
type Schema struct {
	Type       string                            `json:"type" yaml:"type" bson:"type"`
	Properties map[string]map[string]interface{} `json:"properties,omitempty" yaml:"properties,omitempty" bson:"properties,omitempty"`
}

// MockMapping represents a mapping of mocks for schema operations
type MockMapping struct {
	Service   string     `json:"service" yaml:"service" bson:"service"`
	TestSetID string     `json:"testSetId" yaml:"testSetId" bson:"test_set_id"`
	Mocks     []*OpenAPI `json:"mocks" yaml:"mocks" bson:"mocks"`
}

// TestCase represents a test case
type TestCase struct {
	ID       string              `json:"id" yaml:"id" bson:"_id"`
	Name     string              `json:"name" yaml:"name" bson:"name"`
	Created  time.Time           `json:"created" yaml:"created" bson:"created"`
	Updated  time.Time           `json:"updated" yaml:"updated" bson:"updated"`
	Captured time.Time           `json:"captured" yaml:"captured" bson:"captured"`
	CID      string              `json:"cid" yaml:"cid" bson:"cid"`
	AppID    string              `json:"app_id" yaml:"app_id" bson:"app_id"`
	URI      string              `json:"uri" yaml:"uri" bson:"uri"`
	HTTPReq  HTTPRequest         `json:"http_req" yaml:"http_req" bson:"http_req"`
	HTTPResp HTTPResponse        `json:"http_resp" yaml:"http_resp" bson:"http_resp"`
	Deps     []Dependency        `json:"deps" yaml:"deps" bson:"deps"`
	AllKeys  map[string][]string `json:"all_keys" yaml:"all_keys" bson:"all_keys"`
	Anchors  map[string][]string `json:"anchors" yaml:"anchors" bson:"anchors"`
	Noise    []string            `json:"noise" yaml:"noise" bson:"noise"`
}

// Dependency represents a dependency in test cases
type Dependency struct {
	Name string      `json:"name" yaml:"name" bson:"name"`
	Type string      `json:"type" yaml:"type" bson:"type"`
	Meta interface{} `json:"meta" yaml:"meta" bson:"meta"`
}

// SchemaInfo represents schema information for validation
type SchemaInfo struct {
	Service   string  `json:"service" yaml:"service" bson:"service"`
	TestSetID string  `json:"testSetId" yaml:"testSetId" bson:"test_set_id"`
	Name      string  `json:"name" yaml:"name" bson:"name"`
	Score     float64 `json:"score" yaml:"score" bson:"score"`
	Data      OpenAPI `json:"data" yaml:"data" bson:"data"`
}

// Summary represents a validation summary
type Summary struct {
	ServicesSummary []ServiceSummary `json:"servicesSummary" yaml:"servicesSummary" bson:"services_summary"`
}

// ServiceSummary represents a summary for a specific service
type ServiceSummary struct {
	Service     string            `json:"service" yaml:"service" bson:"service"`
	PassedCount int               `json:"passedCount" yaml:"passedCount" bson:"passed_count"`
	FailedCount int               `json:"failedCount" yaml:"failedCount" bson:"failed_count"`
	MissedCount int               `json:"missedCount" yaml:"missedCount" bson:"missed_count"`
	TestSets    map[string]Status `json:"testSets" yaml:"testSets" bson:"test_sets"`
}

// Status represents the status of tests
type Status struct {
	Passed []string `json:"passed" yaml:"passed" bson:"passed"`
	Failed []string `json:"failed" yaml:"failed" bson:"failed"`
	Missed []string `json:"missed" yaml:"missed" bson:"missed"`
}

// Mode constants for validation modes
const (
	IdentifyMode = 0
	CompareMode  = 1
)

// Driven mode constants
const (
	ConsumerMode = "consumer"
	ProviderMode = "provider"
)
