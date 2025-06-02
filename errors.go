package typeid

import "fmt"

// ErrValidation is a sentinel error for validation failures.
// Use errors.Is(err, ErrValidation) to check if an error is a validation error.
var ErrValidation error = &validationError{}

// validationError represents errors that occur during TypeID validation
type validationError struct {
	Message string
	Cause   error // Optional wrapped error (e.g., from base32)
}

// Error implements the error interface
func (e *validationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("typeid: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("typeid: %s", e.Message)
}

// Unwrap returns the underlying cause if present
func (e *validationError) Unwrap() error {
	return e.Cause
}

// Is implements error matching and returns true for any validationError
func (e *validationError) Is(target error) bool {
	return target == ErrValidation
}
