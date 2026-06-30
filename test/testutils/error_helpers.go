package testutils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Error handling guidelines for tests:
// - Use require.NoError for setup operations that must succeed for the test to continue
// - Use assert.NoError for test operations where we want to see all failures
// - Use assert.Error when we expect an error and want to continue testing
// - Always include descriptive error messages for better debugging

// SetupError handles errors that occur during test setup
// Uses require.NoError to fail fast if setup fails.
func SetupError(t testing.TB, err error, operation string) {
	t.Helper()
	require.NoError(t, err, "Setup failed: %s", operation)
}

// TestError handles errors in the main test operations
// Uses assert.NoError to allow other assertions to run.
func TestError(t testing.TB, err error, operation string) {
	t.Helper()
	assert.NoError(t, err, "Test operation failed: %s", operation)
}

// ExpectError verifies that an error occurred as expected.
func ExpectError(t testing.TB, err error, operation string) {
	t.Helper()
	assert.Error(t, err, "Expected error for operation: %s", operation)
}

// ExpectErrorWithMessage verifies error occurred and checks the error message.
func ExpectErrorWithMessage(t testing.TB, err error, expectedMessage, operation string) {
	t.Helper()
	if assert.Error(t, err, "Expected error for operation: %s", operation) {
		assert.Contains(t, err.Error(), expectedMessage,
			"Error message should contain expected text for operation: %s", operation)
	}
}

// ExpectErrorType verifies that a specific type of error occurred.
func ExpectErrorType(t testing.TB, err error, errorType string, operation string) {
	t.Helper()
	if assert.Error(t, err, "Expected %s error for operation: %s", errorType, operation) {
		errorMsg := strings.ToLower(err.Error())
		errorTypeCheck := strings.ToLower(errorType)
		assert.Contains(
			t,
			errorMsg,
			errorTypeCheck,
			"Error should be of type '%s' for operation: %s. Got: %s",
			errorType,
			operation,
			err.Error(),
		)
	}
}

// DatabaseError handles database operation errors with context.
func DatabaseError(t testing.TB, err error, operation, tableName string) {
	t.Helper()
	if err != nil {
		t.Errorf("Database operation failed: %s on table %s. Error: %v", operation, tableName, err)
	}
}

// ServiceError handles service layer errors with context.
func ServiceError(t testing.TB, err error, serviceName, method string) {
	t.Helper()
	assert.NoError(t, err, "Service %s.%s failed", serviceName, method)
}

// RepositoryError handles repository layer errors with context.
func RepositoryError(t testing.TB, err error, repoName, method string) {
	t.Helper()
	assert.NoError(t, err, "Repository %s.%s failed", repoName, method)
}

// ValidationError expects a validation error with specific field.
func ValidationError(t testing.TB, err error, fieldName string) {
	t.Helper()
	if assert.Error(t, err, "Expected validation error for field: %s", fieldName) {
		errorMsg := strings.ToLower(err.Error())
		fieldCheck := strings.ToLower(fieldName)
		assert.Contains(t, errorMsg, fieldCheck,
			"Validation error should mention field '%s'. Got: %s", fieldName, err.Error())
	}
}

// AuthenticationError expects an authentication-related error.
func AuthenticationError(t testing.TB, err error, context string) {
	t.Helper()
	if assert.Error(t, err, "Expected authentication error for: %s", context) {
		errorMsg := strings.ToLower(err.Error())
		authKeywords := []string{"auth", "unauthorized", "forbidden", "permission", "access"}

		hasAuthKeyword := false
		for _, keyword := range authKeywords {
			if strings.Contains(errorMsg, keyword) {
				hasAuthKeyword = true
				break
			}
		}

		assert.True(t, hasAuthKeyword,
			"Error should be authentication-related for: %s. Got: %s", context, err.Error())
	}
}

// NotFoundError expects a "not found" error.
func NotFoundError(t testing.TB, err error, entityType, identifier string) {
	t.Helper()
	if assert.Error(t, err, "Expected not found error for %s: %s", entityType, identifier) {
		errorMsg := strings.ToLower(err.Error())
		assert.Contains(t, errorMsg, "not found",
			"Error should indicate '%s' not found. Got: %s", entityType, err.Error())
	}
}

// ErrorChain helps test error chains and wrapped errors.
type ErrorChain struct {
	t testing.TB
}

// NewErrorChain creates a new error chain tester.
func NewErrorChain(t testing.TB) *ErrorChain {
	return &ErrorChain{t: t}
}

// HasError checks if an error occurred.
func (ec *ErrorChain) HasError(err error, operation string) *ErrorChain {
	ec.t.Helper()
	assert.Error(ec.t, err, "Expected error for operation: %s", operation)
	return ec
}

// NoError checks if no error occurred.
func (ec *ErrorChain) NoError(err error, operation string) *ErrorChain {
	ec.t.Helper()
	assert.NoError(ec.t, err, "Unexpected error for operation: %s", operation)
	return ec
}

// Contains checks if error message contains specific text.
func (ec *ErrorChain) Contains(err error, text string) *ErrorChain {
	ec.t.Helper()
	if err != nil {
		assert.Contains(ec.t, err.Error(), text, "Error message should contain: %s", text)
	}
	return ec
}

// Equals checks if error message equals specific text.
func (ec *ErrorChain) Equals(err error, expectedError error) *ErrorChain {
	ec.t.Helper()
	if expectedError != nil && err != nil {
		assert.Equal(ec.t, expectedError.Error(), err.Error(), "Error messages should match")
	} else {
		assert.Equal(ec.t, expectedError, err, "Errors should be equal")
	}
	return ec
}

// TestScenario helps organize error testing scenarios.
type TestScenario struct {
	Name        string
	Setup       func() error
	Operation   func() error
	Expectation func(testing.TB, error)
}

// RunErrorScenarios runs multiple error test scenarios.
func RunErrorScenarios(t *testing.T, scenarios []TestScenario) {
	t.Helper()

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			// Setup
			if scenario.Setup != nil {
				setupErr := scenario.Setup()
				SetupError(t, setupErr, fmt.Sprintf("setup for scenario '%s'", scenario.Name))
			}

			// Operation
			err := scenario.Operation()

			// Expectation
			scenario.Expectation(t, err)
		})
	}
}

// Common error patterns for reuse.
var (
	// ExpectNoError is a common expectation function.
	ExpectNoError = func(t testing.TB, err error) {
		assert.NoError(t, err, "Expected no error")
	}

	// ExpectAnyError is a common expectation function.
	ExpectAnyError = func(t testing.TB, err error) {
		assert.Error(t, err, "Expected an error")
	}

	// ExpectValidationError is a common expectation function.
	ExpectValidationError = func(t testing.TB, err error) {
		ValidationError(t, err, "validation")
	}

	// ExpectNotFoundError is a common expectation function.
	ExpectNotFoundError = func(t testing.TB, err error) {
		NotFoundError(t, err, "resource", "unknown")
	}
)
