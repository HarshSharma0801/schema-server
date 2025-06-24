package models

type HTTPDoc struct {
	Name    string  `json:"name" yaml:"name" bson:"name"`
	Version string  `json:"version" yaml:"version" bson:"version"`
	Kind    string  `json:"kind" yaml:"kind" bson:"kind"`
	Spec    APISpec `json:"spec" yaml:"spec" bson:"spec"`
}

type APISpec struct {
	Request  HTTPRequest  `json:"request" yaml:"request" bson:"request"`
	Response HTTPResponse `json:"response" yaml:"response" bson:"response"`
}

type HTTPRequest struct {
	Method string            `json:"method" yaml:"method" bson:"method"`
	URL    string            `json:"url" yaml:"url" bson:"url"`
	Header map[string]string `json:"header" yaml:"header" bson:"header"`
	Body   string            `json:"body" yaml:"body" bson:"body"`
}

type HTTPResponse struct {
	StatusCode    int               `json:"statusCode" yaml:"statusCode" bson:"statusCode"`
	StatusMessage string            `json:"statusMessage" yaml:"statusMessage" bson:"statusMessage"`
	Header        map[string]string `json:"header" yaml:"header" bson:"header"`
	Body          string            `json:"body" yaml:"body" bson:"body"`
}

type OpenAPI struct {
	OpenAPI    string                 `json:"openapi" yaml:"openapi" bson:"openapi"`
	Info       Info                   `json:"info" yaml:"info" bson:"info"`
	Servers    []map[string]string    `json:"servers,omitempty" yaml:"servers,omitempty" bson:"servers,omitempty"`
	Paths      map[string]PathItem    `json:"paths" yaml:"paths" bson:"paths"`
	Components map[string]interface{} `json:"components,omitempty" yaml:"components,omitempty" bson:"components,omitempty"`
}

type Info struct {
	Title       string `json:"title" yaml:"title" bson:"title"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" bson:"description,omitempty"`
	Version     string `json:"version" yaml:"version" bson:"version"`
}

type PathItem struct {
	Get    *Operation `json:"get,omitempty" yaml:"get,omitempty" bson:"get,omitempty"`
	Post   *Operation `json:"post,omitempty" yaml:"post,omitempty" bson:"post,omitempty"`
	Put    *Operation `json:"put,omitempty" yaml:"put,omitempty" bson:"put,omitempty"`
	Patch  *Operation `json:"patch,omitempty" yaml:"patch,omitempty" bson:"patch,omitempty"`
	Delete *Operation `json:"delete,omitempty" yaml:"delete,omitempty" bson:"delete,omitempty"`
}

type Operation struct {
	Summary     string                  `json:"summary,omitempty" yaml:"summary,omitempty" bson:"summary,omitempty"`
	Description string                  `json:"description,omitempty" yaml:"description,omitempty" bson:"description,omitempty"`
	OperationID string                  `json:"operationId,omitempty" yaml:"operationId,omitempty" bson:"operationId,omitempty"`
	Parameters  []Parameter             `json:"parameters,omitempty" yaml:"parameters,omitempty" bson:"parameters,omitempty"`
	RequestBody *RequestBody            `json:"requestBody,omitempty" yaml:"requestBody,omitempty" bson:"requestBody,omitempty"`
	Responses   map[string]ResponseItem `json:"responses" yaml:"responses" bson:"responses"`
}

type Parameter struct {
	Name        string      `json:"name" yaml:"name" bson:"name"`
	In          string      `json:"in" yaml:"in" bson:"in"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty" bson:"description,omitempty"`
	Required    bool        `json:"required,omitempty" yaml:"required,omitempty" bson:"required,omitempty"`
	Schema      ParamSchema `json:"schema" yaml:"schema" bson:"schema"`
	Example     interface{} `json:"example,omitempty" yaml:"example,omitempty" bson:"example,omitempty"`
}

type RequestBody struct {
	Content map[string]MediaType `json:"content" yaml:"content" bson:"content"`
}

type ResponseItem struct {
	Description string               `json:"description" yaml:"description" bson:"description"`
	Headers     map[string]Header    `json:"headers,omitempty" yaml:"headers,omitempty" bson:"headers,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty" bson:"content,omitempty"`
}

type Header struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty" bson:"description,omitempty"`
	Schema      Schema `json:"schema" yaml:"schema" bson:"schema"`
}

type MediaType struct {
	Schema  Schema      `json:"schema" yaml:"schema" bson:"schema"`
	Example interface{} `json:"example,omitempty" yaml:"example,omitempty" bson:"example,omitempty"`
}

type Schema struct {
	Type       string                            `json:"type" yaml:"type" bson:"type"`
	Properties map[string]map[string]interface{} `json:"properties,omitempty" yaml:"properties,omitempty" bson:"properties,omitempty"`
}

type ParamSchema struct {
	Type string `json:"type" yaml:"type" bson:"type"`
}
