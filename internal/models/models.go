package models

import "time"

// HTTPDoc represents an HTTP document for schema operations
type HTTPDoc struct {
	Name    string  `json:"name" yaml:"name"`
	Version string  `json:"version" yaml:"version"`
	Kind    string  `json:"kind" yaml:"kind"`
	Spec    APISpec `json:"spec" yaml:"spec"`
}

// APISpec represents the specification of an API
type APISpec struct {
	Request  HTTPRequest  `json:"request" yaml:"request"`
	Response HTTPResponse `json:"response" yaml:"response"`
}

// HTTPRequest represents an HTTP request
type HTTPRequest struct {
	Method string            `json:"method" yaml:"method"`
	URL    string            `json:"url" yaml:"url"`
	Header map[string]string `json:"header" yaml:"header"`
	Body   string            `json:"body" yaml:"body"`
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	StatusCode    int               `json:"statusCode" yaml:"statusCode"`
	StatusMessage string            `json:"statusMessage" yaml:"statusMessage"`
	Header        map[string]string `json:"header" yaml:"header"`
	Body          string            `json:"body" yaml:"body"`
}

// OpenAPI represents an OpenAPI specification
type OpenAPI struct {
	OpenAPI    string                 `json:"openapi" yaml:"openapi"`
	Info       Info                   `json:"info" yaml:"info"`
	Servers    []map[string]string    `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths      map[string]PathItem    `json:"paths" yaml:"paths"`
	Components map[string]interface{} `json:"components,omitempty" yaml:"components,omitempty"`
}

// Info represents the info section of OpenAPI
type Info struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string `json:"version" yaml:"version"`
}

// PathItem represents a path item in OpenAPI
type PathItem struct {
	Get    *Operation `json:"get,omitempty" yaml:"get,omitempty"`
	Post   *Operation `json:"post,omitempty" yaml:"post,omitempty"`
	Put    *Operation `json:"put,omitempty" yaml:"put,omitempty"`
	Patch  *Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
	Delete *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
}

// Operation represents an operation in OpenAPI
type Operation struct {
	Summary     string                  `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string                  `json:"description,omitempty" yaml:"description,omitempty"`
	OperationID string                  `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters  []Parameter             `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody *RequestBody            `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   map[string]ResponseItem `json:"responses" yaml:"responses"`
}

// Parameter represents a parameter in OpenAPI
type Parameter struct {
	Name     string      `json:"name" yaml:"name"`
	In       string      `json:"in" yaml:"in"`
	Required bool        `json:"required,omitempty" yaml:"required,omitempty"`
	Schema   ParamSchema `json:"schema" yaml:"schema"`
	Example  interface{} `json:"example,omitempty" yaml:"example,omitempty"`
}

// ParamSchema represents a parameter schema
type ParamSchema struct {
	Type string `json:"type" yaml:"type"`
}

// RequestBody represents a request body in OpenAPI
type RequestBody struct {
	Content map[string]MediaType `json:"content" yaml:"content"`
}

// ResponseItem represents a response item in OpenAPI
type ResponseItem struct {
	Description string               `json:"description" yaml:"description"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

// MediaType represents a media type in OpenAPI
type MediaType struct {
	Schema  Schema      `json:"schema" yaml:"schema"`
	Example interface{} `json:"example,omitempty" yaml:"example,omitempty"`
}

// Schema represents a schema in OpenAPI
type Schema struct {
	Type       string                            `json:"type" yaml:"type"`
	Properties map[string]map[string]interface{} `json:"properties,omitempty" yaml:"properties,omitempty"`
}

// MockMapping represents a mapping of mocks for schema operations
type MockMapping struct {
	Service   string     `json:"service" yaml:"service"`
	TestSetID string     `json:"testSetId" yaml:"testSetId"`
	Mocks     []*OpenAPI `json:"mocks" yaml:"mocks"`
}

// TestCase represents a test case
type TestCase struct {
	ID       string              `json:"id" yaml:"id"`
	Name     string              `json:"name" yaml:"name"`
	Created  time.Time           `json:"created" yaml:"created"`
	Updated  time.Time           `json:"updated" yaml:"updated"`
	Captured time.Time           `json:"captured" yaml:"captured"`
	CID      string              `json:"cid" yaml:"cid"`
	AppID    string              `json:"app_id" yaml:"app_id"`
	URI      string              `json:"uri" yaml:"uri"`
	HTTPReq  HTTPRequest         `json:"http_req" yaml:"http_req"`
	HTTPResp HTTPResponse        `json:"http_resp" yaml:"http_resp"`
	Deps     []Dependency        `json:"deps" yaml:"deps"`
	AllKeys  map[string][]string `json:"all_keys" yaml:"all_keys"`
	Anchors  map[string][]string `json:"anchors" yaml:"anchors"`
	Noise    []string            `json:"noise" yaml:"noise"`
}

// Dependency represents a dependency in test cases
type Dependency struct {
	Name string      `json:"name" yaml:"name"`
	Type string      `json:"type" yaml:"type"`
	Meta interface{} `json:"meta" yaml:"meta"`
}

// SchemaInfo represents schema information for validation
type SchemaInfo struct {
	Service   string  `json:"service" yaml:"service"`
	TestSetID string  `json:"testSetId" yaml:"testSetId"`
	Name      string  `json:"name" yaml:"name"`
	Score     float64 `json:"score" yaml:"score"`
	Data      OpenAPI `json:"data" yaml:"data"`
}

// Summary represents a validation summary
type Summary struct {
	ServicesSummary []ServiceSummary `json:"servicesSummary" yaml:"servicesSummary"`
}

// ServiceSummary represents a summary for a specific service
type ServiceSummary struct {
	Service     string            `json:"service" yaml:"service"`
	PassedCount int               `json:"passedCount" yaml:"passedCount"`
	FailedCount int               `json:"failedCount" yaml:"failedCount"`
	MissedCount int               `json:"missedCount" yaml:"missedCount"`
	TestSets    map[string]Status `json:"testSets" yaml:"testSets"`
}

// Status represents the status of tests
type Status struct {
	Passed []string `json:"passed" yaml:"passed"`
	Failed []string `json:"failed" yaml:"failed"`
	Missed []string `json:"missed" yaml:"missed"`
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
