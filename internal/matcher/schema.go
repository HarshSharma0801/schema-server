// Package matcher provides schema matching functionality based on Keploy's matcher logic
package matcher

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"schema-server/internal/models"

	"go.uber.org/zap"
)

// SchemaMatchMode represents the mode for schema matching
type SchemaMatchMode int

const (
	// CompareMode is for detailed comparison with diff output
	CompareMode SchemaMatchMode = iota
	// IdentifyMode is for identification and scoring
	IdentifyMode
)

const NOTCANDIDATE = -1.0

// ValidatedJSON holds validated JSON data for comparison
type ValidatedJSON struct {
	expected    interface{}
	actual      interface{}
	isIdentical bool
}

func (v *ValidatedJSON) IsIdentical() bool {
	return v.isIdentical
}

func (v *ValidatedJSON) Expected() interface{} {
	return v.expected
}

func (v *ValidatedJSON) Actual() interface{} {
	return v.actual
}

// JSONComparisonResult holds the result of JSON comparison
type JSONComparisonResult struct {
	matches     bool
	isExact     bool
	differences []string
}

func (v *JSONComparisonResult) IsExact() bool {
	return v.isExact
}

func (v *JSONComparisonResult) Matches() bool {
	return v.matches
}

func (v *JSONComparisonResult) Differences() []string {
	return v.differences
}

// DiffsPrinter is a simplified version for logging diffs
type DiffsPrinter struct {
	testCase              string
	hasarrayIndexMismatch bool
	text                  string
}

func NewDiffsPrinter(testCase string) DiffsPrinter {
	return DiffsPrinter{testCase: testCase}
}

func (d *DiffsPrinter) SetHasarrayIndexMismatch(has bool) {
	d.hasarrayIndexMismatch = has
}

func (d *DiffsPrinter) PushTypeDiff(exp, act string) {
	// Simplified implementation
}

func (d *DiffsPrinter) PushFooterDiff(key string) {
	d.hasarrayIndexMismatch = true
	d.text = key
}

func (d *DiffsPrinter) PushBodyDiff(exp, act string, noise map[string][]string) {
	// Simplified implementation
}

func (d *DiffsPrinter) RenderAppender() error {
	// Simplified implementation - just return nil
	return nil
}

// Match compares two OpenAPI schemas and returns a score and match result
// This is the main function ported from Keploy's matcher logic
func Match(mock, test models.OpenAPI, testSetID string, mockSetID string, logger *zap.Logger, mode SchemaMatchMode) (float64, bool, error) {
	// Ignore info/version/description for matching
	mockCopy := mock
	testCopy := test
	mockCopy.Info = models.Info{}
	testCopy.Info = models.Info{}

	// Only compare paths and components.schemas semantically
	mockSchemas := getSchemasMap(mockCopy.Components)
	testSchemas := getSchemasMap(testCopy.Components)

	if semanticEqual(mockCopy.Paths, testCopy.Paths) && semanticEqual(mockSchemas, testSchemas) {
		return 0, true, nil
	}
	return -1, false, nil
}

// getSchemasMap extracts the schemas map from components (if present)
func getSchemasMap(components map[string]interface{}) map[string]interface{} {
	if components == nil {
		return nil
	}
	if schemas, ok := components["schemas"]; ok {
		if m, ok := schemas.(map[string]interface{}); ok {
			return m
		}
	}
	return nil
}

// semanticEqual recursively compares two values, ignoring order in maps/arrays and allowing extra optional fields
func semanticEqual(a, b interface{}) bool {
	switch aVal := a.(type) {
	case map[string]interface{}:
		bVal, ok := b.(map[string]interface{})
		if !ok {
			return false
		}
		for k, v := range aVal {
			if !semanticEqual(v, bVal[k]) {
				return false
			}
		}
		for k, v := range bVal {
			if _, ok := aVal[k]; !ok {
				// extra field in b, allowed
				continue
			}
			if !semanticEqual(aVal[k], v) {
				return false
			}
		}
		return true
	case []interface{}:
		bVal, ok := b.([]interface{})
		if !ok {
			return false
		}
		if len(aVal) != len(bVal) {
			return false
		}
		// Compare ignoring order
		used := make([]bool, len(bVal))
		for _, av := range aVal {
			found := false
			for i, bv := range bVal {
				if used[i] {
					continue
				}
				if semanticEqual(av, bv) {
					used[i] = true
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(a, b)
	}
}

func compareOperationTypes(mockOperationType, testOperationType string) (bool, error) {
	pass := true
	if mockOperationType != testOperationType {
		pass = false
		return pass, nil
	}
	return pass, nil
}

func compareRequestBodies(mockOperation, testOperation *models.Operation, logDiffs DiffsPrinter, logger *zap.Logger, testName, mockName, testSetID, mockSetID string) (bool, error) {
	pass := false
	var score float64

	mockRequestBodyStr, testRequestBodyStr, err := MarshalRequestBodies(mockOperation, testOperation)
	if err != nil {
		return false, err
	}

	validatedJSON, err := ValidateAndMarshalJSON(logger, &mockRequestBodyStr, &testRequestBodyStr)
	if err != nil {
		return false, err
	}

	if validatedJSON.IsIdentical() {
		if score, pass, err = handleJSONDiff(validatedJSON, logDiffs, logger, testName, mockName, testSetID, mockSetID, mockRequestBodyStr, testRequestBodyStr, "request", 0); err != nil {
			return false, err
		}
		if score == NOTCANDIDATE {
			return false, nil
		}
	} else {
		pass = false
		return pass, nil
	}
	return pass, nil
}

func compareParameters(mockParameters, testParameters []models.Parameter) (bool, error) {
	pass := true

	for _, mockParam := range mockParameters {
		if mockParam.In == "header" {
			continue
		}
		found := false
		for _, testParam := range testParameters {
			if mockParam.Name == testParam.Name && mockParam.In == testParam.In {
				found = true
				if mockParam.Schema.Type != testParam.Schema.Type {
					pass = false
					return pass, nil
				}
			}
		}
		if !found {
			pass = false
			return pass, nil
		}
	}

	return pass, nil
}

func compareResponseBodies(status string, mockOperation, testOperation *models.Operation, logDiffs DiffsPrinter, logger *zap.Logger, testName, mockName, testSetID, mockSetID string, mode SchemaMatchMode) (float64, bool, bool, error) {
	pass := true
	overallScore := 0.0
	matched := false
	differencesCount := 0.0

	if _, ok := testOperation.Responses[status]; ok {
		mockResponseBodyStr, testResponseBodyStr, err := MarshalResponseBodies(status, mockOperation, testOperation)
		if err != nil {
			return differencesCount, false, false, err
		}

		if mockOperation.Responses[status].Content["application/json"].Schema.Properties != nil {
			overallScore = float64(len(mockOperation.Responses[status].Content["application/json"].Schema.Properties))
		}

		validatedJSON, err := ValidateAndMarshalJSON(logger, &mockResponseBodyStr, &testResponseBodyStr)
		if err != nil {
			return differencesCount, false, false, err
		}

		if validatedJSON.IsIdentical() {
			switch mode {
			case CompareMode:
				if _, matched, err = handleJSONDiff(validatedJSON, logDiffs, logger, testName, mockName, testSetID, mockSetID, mockResponseBodyStr, testResponseBodyStr, "response", mode); err != nil {
					return differencesCount, false, false, err
				}
			case IdentifyMode:
				differencesCount, err = calculateSimilarityScore(mockOperation, testOperation, status)
				if err != nil {
					return differencesCount, false, false, err
				}
			}
		} else {
			differencesCount = overallScore

			if mode == CompareMode {
				logDiffs.PushTypeDiff(fmt.Sprint(reflect.TypeOf(validatedJSON.Expected())), fmt.Sprint(reflect.TypeOf(validatedJSON.Actual())))
				logger.Warn("Contract Check failed",
					zap.String("test", testName),
					zap.String("testSetID", testSetID),
					zap.String("mock", mockName),
					zap.String("mockSetID", mockSetID))
			}
		}
	} else {
		pass = false
		differencesCount = -1
	}

	if overallScore > 0 {
		return differencesCount / overallScore, pass, matched, nil
	}
	return differencesCount, pass, matched, nil
}

func calculateSimilarityScore(mockOperation, testOperation *models.Operation, status string) (float64, error) {
	testContent := testOperation.Responses[status].Content["application/json"]
	mockContent := mockOperation.Responses[status].Content["application/json"]

	if testContent.Schema.Properties == nil || mockContent.Schema.Properties == nil {
		return 0.0, nil
	}

	testParameters := testContent.Schema.Properties
	mockParameters := mockContent.Schema.Properties
	score := 0.0

	for key, testParam := range testParameters {
		if mockParam, ok := mockParameters[key]; ok {
			if testParam["type"] == mockParam["type"] {
				score++
			}
		}
	}
	return score, nil
}

func handleJSONDiff(validatedJSON ValidatedJSON, logDiffs DiffsPrinter, logger *zap.Logger, testName, mockName, testSetID, mockSetID, mockBodyStr, testBodyStr, diffType string, mode SchemaMatchMode) (float64, bool, error) {
	pass := true
	differencesCount := 0.0

	jsonComparisonResult, err := JSONDiffWithNoiseControl(validatedJSON, nil, false)
	if err != nil {
		return differencesCount, false, err
	}

	if !jsonComparisonResult.IsExact() {
		pass = false
		if json.Valid([]byte(mockBodyStr)) {
			// Simple diff count - in real implementation this would use jsondiff
			if mockBodyStr != testBodyStr {
				differencesCount = 1.0
			}

			if diffType == "request" && differencesCount > 1 {
				return -1.0, false, nil
			}

			if diffType == "response" {
				if jsonComparisonResult.Matches() {
					logDiffs.SetHasarrayIndexMismatch(true)
					logDiffs.PushFooterDiff(strings.Join(jsonComparisonResult.Differences(), ", "))
				}
				logDiffs.PushBodyDiff(mockBodyStr, testBodyStr, nil)
			}
		}

		if diffType == "response" && mode == CompareMode {
			if err := logDiffs.RenderAppender(); err != nil {
				return differencesCount, false, err
			}
		}
	}
	return differencesCount, pass, nil
}
