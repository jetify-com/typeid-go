package typeid

import (
	"fmt"

	"go.jetify.com/typeid/v2/base32"
)

func validatePrefix(prefix string) error {
	if len(prefix) > 63 {
		return &validationError{
			Message: fmt.Sprintf("prefix length must be <= 63, got %d for %q", len(prefix), prefix),
		}
	}

	if len(prefix) > 0 && prefix[0] == '_' {
		return &validationError{
			Message: fmt.Sprintf("prefix cannot start with underscore, got %q", prefix),
		}
	}

	if len(prefix) > 0 && prefix[len(prefix)-1] == '_' {
		return &validationError{
			Message: fmt.Sprintf("prefix cannot end with underscore, got %q", prefix),
		}
	}

	// Ensure that the prefix only has lowercase ASCII characters
	for _, c := range prefix {
		if (c < 'a' || c > 'z') && c != '_' {
			return &validationError{
				Message: fmt.Sprintf("prefix must contain only [a-z_], found %q in %q", c, prefix),
			}
		}
	}

	return nil
}

func validateSuffix(suffix string) error {
	if len(suffix) != 26 {
		return &validationError{
			Message: fmt.Sprintf("suffix length must be 26, got %d", len(suffix)),
		}
	}

	if suffix[0] > '7' {
		return &validationError{
			Message: fmt.Sprintf("suffix must start with 0-7, got %q", suffix[0]),
		}
	}
	// Validate the suffix using zero-allocation validation
	if err := base32.ValidateString(suffix); err != nil {
		return &validationError{
			Message: "invalid suffix encoding",
			Cause:   err,
		}
	}
	return nil
}
