package typeid

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateValidationErrors(t *testing.T) {
	tests := []struct {
		name             string
		prefix           string
		expectValidation bool
	}{
		{
			name:             "prefix too long",
			prefix:           "this_is_a_very_long_prefix_that_exceeds_the_sixty_three_character_limit",
			expectValidation: true,
		},
		{
			name:             "prefix with uppercase",
			prefix:           "PREFIX",
			expectValidation: true,
		},
		{
			name:             "prefix with number",
			prefix:           "prefix123",
			expectValidation: true,
		},
		{
			name:             "prefix with space",
			prefix:           "prefix space",
			expectValidation: true,
		},
		{
			name:             "prefix with special char",
			prefix:           "prefix-dash",
			expectValidation: true,
		},
		{
			name:             "prefix starts with underscore",
			prefix:           "_prefix",
			expectValidation: true,
		},
		{
			name:             "prefix ends with underscore",
			prefix:           "prefix_",
			expectValidation: true,
		},
		{
			name:             "valid prefix",
			prefix:           "valid_prefix",
			expectValidation: false,
		},
		{
			name:             "empty prefix",
			prefix:           "",
			expectValidation: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Generate(tt.prefix)

			if tt.expectValidation {
				require.Error(t, err, "expected validation error")
				assert.True(t, errors.Is(err, ErrValidation), "expected ErrValidation, got %v", err)
			} else {
				assert.NoError(t, err, "expected no error")
			}
		})
	}
}

func TestParseValidationErrors(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectValidation bool
	}{
		{
			name:             "suffix too short",
			input:            "test_invalid",
			expectValidation: true,
		},
		{
			name:             "suffix too long",
			input:            "test_000000000000000000000000000",
			expectValidation: true,
		},
		{
			name:             "suffix starts with 8",
			input:            "test_80000000000000000000000000",
			expectValidation: true,
		},
		{
			name:             "suffix starts with 9",
			input:            "test_90000000000000000000000000",
			expectValidation: true,
		},
		{
			name:             "suffix with invalid base32 char !",
			input:            "test_0000000000000000000000!0",
			expectValidation: true,
		},
		{
			name:             "suffix with invalid base32 char @",
			input:            "test_0123456789012345678901234@",
			expectValidation: true,
		},
		{
			name:             "suffix with uppercase",
			input:            "test_0123456789012345678901234A",
			expectValidation: true,
		},
		{
			name:             "empty prefix with separator",
			input:            "_00000000000000000000000000",
			expectValidation: true,
		},
		{
			name:             "prefix with invalid char in parse",
			input:            "PREFIX_00000000000000000000000000",
			expectValidation: true,
		},
		{
			name:             "valid typeid",
			input:            "test_00000000000000000000000000",
			expectValidation: false,
		},
		{
			name:             "valid typeid no prefix",
			input:            "00000000000000000000000000",
			expectValidation: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)

			if tt.expectValidation {
				require.Error(t, err, "expected validation error")
				assert.True(t, errors.Is(err, ErrValidation), "expected ErrValidation, got %v", err)
			} else {
				assert.NoError(t, err, "expected no error")
			}
		})
	}
}

func TestFromUUIDValidationErrors(t *testing.T) {
	tests := []struct {
		name             string
		uuid             string
		prefix           string
		expectValidation bool
	}{
		{
			name:             "invalid uuid format",
			uuid:             "not-a-uuid",
			prefix:           "test",
			expectValidation: true, // UUID parsing error is now a validation error
		},
		{
			name:             "valid uuid with invalid prefix",
			uuid:             "00000000-0000-0000-0000-000000000000",
			prefix:           "INVALID",
			expectValidation: true,
		},
		{
			name:             "valid uuid with valid prefix",
			uuid:             "00000000-0000-0000-0000-000000000000",
			prefix:           "test",
			expectValidation: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromUUID(tt.prefix, tt.uuid)

			if err != nil {
				isValidation := errors.Is(err, ErrValidation)
				assert.Equal(t, tt.expectValidation, isValidation,
					"expected validation=%v for error %v", tt.expectValidation, err)
			} else {
				assert.False(t, tt.expectValidation, "expected error but got none")
			}
		})
	}
}

func TestInternalParseValidationErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty suffix",
			input: "test_",
		},
		{
			name:  "invalid prefix",
			input: "TEST_00000000000000000000000000",
		},
		{
			name:  "invalid suffix length",
			input: "test_short",
		},
		{
			name:  "empty prefix with separator",
			input: "_00000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			require.Error(t, err, "expected error")
			assert.True(t, errors.Is(err, ErrValidation), "expected validation error")
		})
	}
}

func TestValidationErrorUnwrap(t *testing.T) {
	t.Run("with cause", func(t *testing.T) {
		// Test that base32 errors are properly wrapped
		// Using uppercase 'A' which is invalid in base32
		_, err := Parse("test_0123456789012345678901234A")
		require.Error(t, err)

		assert.True(t, errors.Is(err, ErrValidation), "expected validation error")

		// Should be able to unwrap to get the base32 error
		var valErr *validationError
		require.True(t, errors.As(err, &valErr), "expected to unwrap validation error")
		require.NotNil(t, valErr.Cause, "expected validation error to have a cause for base32 decoding error")

		// Test Unwrap method directly
		unwrapped := valErr.Unwrap()
		assert.NotNil(t, unwrapped, "Unwrap() should return non-nil error")
		assert.Equal(t, valErr.Cause, unwrapped, "Unwrap() should return the cause")
	})

	t.Run("without cause", func(t *testing.T) {
		err := &validationError{Message: "no cause"}
		assert.Nil(t, err.Unwrap(), "Unwrap() should return nil when there's no cause")
	})
}

func TestValidationErrorMessages(t *testing.T) {
	tests := []struct {
		name string
		err  *validationError
		want string
	}{
		{
			name: "error without cause",
			err: &validationError{
				Message: "test error",
			},
			want: "typeid: test error",
		},
		{
			name: "error with cause",
			err: &validationError{
				Message: "wrapper error",
				Cause:   errors.New("underlying error"),
			},
			want: "typeid: wrapper error: underlying error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidationErrorIs(t *testing.T) {
	tests := []struct {
		name   string
		err    *validationError
		target error
		want   bool
	}{
		{
			name:   "matches ErrValidation sentinel",
			err:    &validationError{Message: "any message"},
			target: ErrValidation,
			want:   true,
		},
		{
			name:   "does not match other errors",
			err:    &validationError{Message: "any message"},
			target: errors.New("some other error"),
			want:   false,
		},
		{
			name:   "nil target",
			err:    &validationError{Message: "any message"},
			target: nil,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Is(tt.target)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestErrValidationSentinel(t *testing.T) {
	// Test that ErrValidation is properly initialized
	assert.NotNil(t, ErrValidation)
	assert.IsType(t, &validationError{}, ErrValidation)

	// Test that errors.Is works with the sentinel
	err1 := &validationError{Message: "test 1"}
	err2 := &validationError{Message: "test 2"}

	assert.True(t, errors.Is(err1, ErrValidation))
	assert.True(t, errors.Is(err2, ErrValidation))

	// Test that different validation errors are equal to the sentinel but not to each other
	assert.False(t, errors.Is(err1, err2))
}

func TestGenerateErrorMessages(t *testing.T) {
	tests := []struct {
		name            string
		prefix          string
		expectedPattern string
	}{
		{
			name:            "prefix too long",
			prefix:          "this_is_a_very_long_prefix_that_exceeds_the_sixty_three_character_limit",
			expectedPattern: "prefix length must be <= 63",
		},
		{
			name:            "prefix with dash",
			prefix:          "has-dash",
			expectedPattern: "prefix must contain only [a-z_]",
		},
		{
			name:            "prefix with uppercase",
			prefix:          "PREFIX",
			expectedPattern: "prefix must contain only [a-z_]",
		},
		{
			name:            "prefix starts with underscore",
			prefix:          "_invalid",
			expectedPattern: "prefix cannot start with underscore",
		},
		{
			name:            "prefix ends with underscore",
			prefix:          "invalid_",
			expectedPattern: "prefix cannot end with underscore",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Generate(tt.prefix)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedPattern,
				"error message should contain expected pattern")
			assert.Contains(t, err.Error(), "typeid:",
				"error message should have typeid prefix")
		})
	}
}

func TestParseErrorMessages(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedPattern string
	}{
		{
			name:            "suffix too short",
			input:           "test_short",
			expectedPattern: "suffix length must be 26",
		},
		{
			name:            "suffix too long",
			input:           "test_000000000000000000000000000",
			expectedPattern: "suffix length must be 26",
		},
		{
			name:            "suffix overflow with 8",
			input:           "test_80000000000000000000000000",
			expectedPattern: "suffix must start with 0-7",
		},
		{
			name:            "suffix overflow with 9",
			input:           "test_90000000000000000000000000",
			expectedPattern: "suffix must start with 0-7",
		},
		{
			name:            "base32 invalid char @",
			input:           "test_0123456789012345678901234@",
			expectedPattern: "invalid suffix encoding",
		},
		{
			name:            "base32 invalid uppercase",
			input:           "test_0123456789012345678901234A",
			expectedPattern: "invalid suffix encoding",
		},
		{
			name:            "empty prefix with separator",
			input:           "_00000000000000000000000000",
			expectedPattern: "prefix cannot be empty when separator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedPattern,
				"error message should contain expected pattern")
			assert.Contains(t, err.Error(), "typeid:",
				"error message should have typeid prefix")
		})
	}
}

func TestFromUUIDErrorMessages(t *testing.T) {
	tests := []struct {
		name            string
		uuid            string
		prefix          string
		expectedPattern string
	}{
		{
			name:            "invalid uuid format",
			uuid:            "not-a-uuid",
			prefix:          "test",
			expectedPattern: "invalid UUID format",
		},
		{
			name:            "invalid uuid with special chars",
			uuid:            "!!!-!!!",
			prefix:          "test",
			expectedPattern: "invalid UUID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromUUID(tt.prefix, tt.uuid)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedPattern,
				"error message should contain expected pattern")
			assert.Contains(t, err.Error(), "typeid:",
				"error message should have typeid prefix")
		})
	}
}
