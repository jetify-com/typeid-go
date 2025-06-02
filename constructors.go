package typeid

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
	"go.jetify.com/typeid/base32"
)

// Generate returns a new TypeID with the given prefix and a random suffix.
// If you want to create an id without a prefix, pass an empty string.
func Generate(prefix string) (TypeID, error) {
	// Validate prefix early
	if err := validatePrefix(prefix); err != nil {
		return zeroID, err
	}

	// Generate new UUID v7
	uid, err := uuid.NewV7()
	if err != nil {
		return zeroID, err
	}

	// Use stack buffer for base32 encoding to avoid allocation
	var suffixBuf [26]byte
	base32.Encode(suffixBuf[:], [16]byte(uid))

	// Build TypeID using helper function
	return newTypeID(prefix, suffixBuf), nil
}

// MustGenerate returns a new TypeID with the given prefix and a random suffix.
// It panics if the prefix is invalid. Use Generate() if you need error handling.
func MustGenerate(prefix string) TypeID {
	tid, err := Generate(prefix)
	if err != nil {
		panic(err)
	}
	return tid
}

// Parse parses a TypeID from a string of the form <prefix>_<suffix>
func Parse(s string) (TypeID, error) {
	prefix, suffix, err := split(s)
	if err != nil {
		return zeroID, err
	}

	// Validate prefix
	if err := validatePrefix(prefix); err != nil {
		return zeroID, err
	}

	// Build TypeID from string parts
	if suffix == "" {
		return zeroID, &validationError{
			Message: "suffix cannot be empty",
		}
	}

	// Validate suffix
	if err := validateSuffix(suffix); err != nil {
		return zeroID, err
	}

	// Handle zero suffix case - empty TypeID should be functionally equivalent
	if prefix == "" && suffix == ZeroSuffix {
		// Return zero TypeID for compatibility with tests
		return zeroID, nil
	}

	// Build TypeID efficiently
	var tid TypeID
	if prefix == "" {
		tid.value = suffix
		tid.prefixLen = 0
	} else {
		tid.value = prefix + "_" + suffix
		tid.prefixLen = uint8(len(prefix))
	}
	return tid, nil
}

func split(id string) (string, string, error) {
	index := strings.LastIndex(id, "_")
	if index == -1 {
		return "", id, nil
	}

	prefix := id[:index]
	suffix := id[index+1:]
	if prefix == "" {
		return "", "", &validationError{
			Message: "prefix cannot be empty when separator \"_\" is present",
		}
	}
	return prefix, suffix, nil
}

// newTypeID constructs a TypeID from a prefix and 26-byte base32 suffix
// with minimal allocations. The suffix must be exactly 26 bytes.
func newTypeID(prefix string, suffixBuf [26]byte) TypeID {
	var tid TypeID
	if prefix == "" {
		tid.value = string(suffixBuf[:])
		tid.prefixLen = 0
	} else {
		var builder strings.Builder
		totalLen := len(prefix) + 1 + len(suffixBuf)
		builder.Grow(totalLen) // Pre-allocate capacity to avoid reallocations
		builder.WriteString(prefix)
		builder.WriteByte('_')
		builder.Write(suffixBuf[:])
		tid.value = builder.String()
		tid.prefixLen = uint8(len(prefix))
	}
	return tid
}

// FromUUID encodes the given UUID (in hex string form) as a TypeID with the given prefix.
// If you want to create an id without a prefix, pass an empty string for the prefix.
func FromUUID(prefix string, uidStr string) (TypeID, error) {
	// Validate prefix early
	if err := validatePrefix(prefix); err != nil {
		return zeroID, err
	}

	uid, err := uuid.FromString(uidStr)
	if err != nil {
		return zeroID, &validationError{
			Message: fmt.Sprintf("invalid UUID format %q", uidStr),
			Cause:   err,
		}
	}

	// Handle zero UUID case - return canonical zeroID for consistency
	if uid == (uuid.UUID{}) && prefix == "" {
		return zeroID, nil
	}

	// Use stack buffer for base32 encoding to avoid allocation
	var suffixBuf [26]byte
	base32.Encode(suffixBuf[:], [16]byte(uid))

	// No need to validate suffixBuf - base32.Encode() always produces valid 26-byte output

	// Build TypeID using helper function
	return newTypeID(prefix, suffixBuf), nil
}

// FromBytes creates a TypeID from a prefix and 16-byte UUID with zero allocations.
// The bytes must be exactly 16 bytes long (standard UUID byte length).
// If you want to create an id without a prefix, pass an empty string for the prefix.
func FromBytes(prefix string, uidBytes []byte) (TypeID, error) {
	// Validate inputs early
	if len(uidBytes) != 16 {
		return zeroID, &validationError{
			Message: fmt.Sprintf("UUID bytes must be exactly 16 bytes, got %d", len(uidBytes)),
		}
	}

	if err := validatePrefix(prefix); err != nil {
		return zeroID, err
	}

	// Handle zero UUID case - return canonical zeroID for consistency
	if isZeroBytes(uidBytes) && prefix == "" {
		return zeroID, nil
	}

	// Convert to array for base32 encoding (zero allocation conversion)
	var uidArray [16]byte
	copy(uidArray[:], uidBytes)

	// Use stack buffer for base32 encoding to avoid allocation
	var suffixBuf [26]byte
	base32.Encode(suffixBuf[:], uidArray)

	// Build TypeID using helper function
	return newTypeID(prefix, suffixBuf), nil
}

// isZeroBytes efficiently checks if a 16-byte slice is all zeros
func isZeroBytes(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}
