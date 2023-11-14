package typeid

import (
	"fmt"

	"go.jetpack.io/typeid/base32"
)

func validatePrefix[T Subtype](prefix string) error {
	if len(prefix) > 63 {
		return fmt.Errorf("invalid prefix: %s. Prefix length is %d, expected <= 63", prefix, len(prefix))
	}

	// Ensure that the prefix only has lowercase ASCII characters
	for _, c := range prefix {
		if c < 'a' || c > 'z' {
			return fmt.Errorf("invalid prefix: '%s'. Prefix should match [a-z]{0,63}", prefix)
		}
	}

	if !isAnyID[T]() {
		expected := defaultType[T]()
		if expected != prefix {
			return fmt.Errorf("invalid prefix: '%s'. Subtype requires prefix to match '%s'", prefix, expected)
		}
	}

	return nil
}

func validateSuffix(suffix string) error {
	if len(suffix) != 26 {
		return fmt.Errorf("invalid suffix: %s. Suffix length is %d, expected 26", suffix, len(suffix))
	}

	if suffix[0] > '7' {
		return fmt.Errorf("invalid suffix: '%s'. Suffix must start with a 0-7 digit to avoid overflows", suffix)
	}
	// Validate the suffix by decoding it, it must be a valid base32 string
	if _, err := base32.Decode(suffix); err != nil {
		return fmt.Errorf("invalid suffix: %w", err)
	}
	return nil
}
