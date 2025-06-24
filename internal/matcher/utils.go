package matcher

import (
	"encoding/json"
	"reflect"

	"schema-server/internal/models"

	"go.uber.org/zap"
)

// FindOperation finds the first non-nil operation in a PathItem and returns it with its method
func FindOperation(item models.PathItem) (*models.Operation, string) {
	operations := map[string]*models.Operation{
		"GET":    item.Get,
		"POST":   item.Post,
		"PUT":    item.Put,
		"DELETE": item.Delete,
		"PATCH":  item.Patch,
	}

	for method, operation := range operations {
		if operation != nil {
			return operation, method
		}
	}
	return nil, ""
}

// MarshalRequestBodies marshals request bodies from operations to JSON strings
func MarshalRequestBodies(mockOperation, testOperation *models.Operation) (string, string, error) {
	var mockRequestBody, testRequestBody string

	if mockOperation.RequestBody != nil && mockOperation.RequestBody.Content["application/json"].Example != nil {
		mockBytes, err := json.Marshal(mockOperation.RequestBody.Content["application/json"].Example)
		if err != nil {
			return "", "", err
		}
		mockRequestBody = string(mockBytes)
	}

	if testOperation.RequestBody != nil && testOperation.RequestBody.Content["application/json"].Example != nil {
		testBytes, err := json.Marshal(testOperation.RequestBody.Content["application/json"].Example)
		if err != nil {
			return "", "", err
		}
		testRequestBody = string(testBytes)
	}

	return mockRequestBody, testRequestBody, nil
}

// MarshalResponseBodies marshals response bodies from operations to JSON strings
func MarshalResponseBodies(status string, mockOperation, testOperation *models.Operation) (string, string, error) {
	var mockResponseBody, testResponseBody string

	if mockResp, ok := mockOperation.Responses[status]; ok {
		if mockResp.Content["application/json"].Example != nil {
			mockBytes, err := json.Marshal(mockResp.Content["application/json"].Example)
			if err != nil {
				return "", "", err
			}
			mockResponseBody = string(mockBytes)
		}
	}

	if testResp, ok := testOperation.Responses[status]; ok {
		if testResp.Content["application/json"].Example != nil {
			testBytes, err := json.Marshal(testResp.Content["application/json"].Example)
			if err != nil {
				return "", "", err
			}
			testResponseBody = string(testBytes)
		}
	}

	return mockResponseBody, testResponseBody, nil
}

// ValidateAndMarshalJSON validates and marshals JSON strings for comparison
func ValidateAndMarshalJSON(log *zap.Logger, exp, act *string) (ValidatedJSON, error) {
	var validatedJSON ValidatedJSON
	var expected interface{}
	var actual interface{}
	var err error

	if *exp != "" {
		expected, err = UnmarshallJSON(*exp, log)
		if err != nil {
			return validatedJSON, err
		}
	}
	if *act != "" {
		actual, err = UnmarshallJSON(*act, log)
		if err != nil {
			return validatedJSON, err
		}
	}

	validatedJSON.expected = expected
	validatedJSON.actual = actual

	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		validatedJSON.isIdentical = false
		return validatedJSON, nil
	}

	cleanExp, err := json.Marshal(expected)
	if err != nil {
		return validatedJSON, err
	}
	cleanAct, err := json.Marshal(actual)
	if err != nil {
		return validatedJSON, err
	}

	*exp = string(cleanExp)
	*act = string(cleanAct)
	validatedJSON.isIdentical = true
	return validatedJSON, nil
}

// UnmarshallJSON returns unmarshalled JSON object
func UnmarshallJSON(s string, log *zap.Logger) (interface{}, error) {
	var result interface{}
	if s == "" {
		return nil, nil
	}
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		log.Error("cannot convert json string into json object", zap.String("json", s), zap.Error(err))
		return nil, err
	}
	return result, nil
}

// JSONDiffWithNoiseControl performs JSON comparison with noise control
func JSONDiffWithNoiseControl(validatedJSON ValidatedJSON, noise map[string][]string, ignoreOrdering bool) (JSONComparisonResult, error) {
	var result JSONComparisonResult

	// Simple implementation - in a full implementation this would have more sophisticated diff logic
	expected := validatedJSON.Expected()
	actual := validatedJSON.Actual()

	// Check if they are deeply equal
	if reflect.DeepEqual(expected, actual) {
		result.matches = true
		result.isExact = true
		result.differences = []string{}
	} else {
		result.matches = false
		result.isExact = false
		result.differences = []string{"structures differ"}

		// If ignoring ordering, try to see if they match when treated as unordered
		if ignoreOrdering {
			// This is a simplified check - in reality this would be more complex
			expectedJSON, _ := json.Marshal(expected)
			actualJSON, _ := json.Marshal(actual)

			if string(expectedJSON) == string(actualJSON) {
				result.matches = true
			}
		}
	}

	return result, nil
}
