package typeid

import (
	"github.com/gofrs/uuid/v5"
	"go.jetify.com/typeid/base32"
)

// TypeID is a unique identifier with a given type as defined by the TypeID spec
type TypeID struct {
	value     string // The complete "prefix_suffix" string, where suffix is the base32 representation.
	prefixLen uint8  // Length of prefix for extracting parts
}

// Prefix returns the type prefix of the TypeID
func (tid TypeID) Prefix() string {
	if tid.prefixLen == 0 {
		return ""
	}
	return tid.value[:tid.prefixLen]
}

const ZeroSuffix = "00000000000000000000000000"

// zeroID is a zero-value TypeID with empty prefix and zero UUID.
var zeroID = TypeID{}

// Suffix returns the suffix of the TypeID in it's canonical base32 representation.
func (tid TypeID) Suffix() string {
	if tid.prefixLen == 0 {
		// No prefix, entire value is the suffix
		if tid.value == "" {
			return ZeroSuffix
		}
		return tid.value
	}
	// Has prefix, suffix starts after prefix + underscore
	if len(tid.value) <= int(tid.prefixLen)+1 {
		return ZeroSuffix
	}
	suffix := tid.value[tid.prefixLen+1:]
	if suffix == "" {
		return ZeroSuffix
	}
	return suffix
}

// String returns the TypeID in it's canonical string representation of the form:
// <prefix>_<suffix> where <suffix> is the canonical base32 representation of the UUID
func (tid TypeID) String() string {
	if tid.value == "" {
		return ZeroSuffix
	}
	return tid.value
}

// Bytes decodes the TypeID's suffix as a UUID and returns it's bytes
func (tid TypeID) Bytes() []byte {
	suffix := tid.Suffix()
	var dst [16]byte
	_, err := base32.Decode(dst[:], []byte(suffix))
	// Decode only fails if the suffix cannot be decoded for one of two reasons:
	// 1. The suffix is not 26 characters long
	// 2. The suffix contains characters that are not in the base32 alphabet
	// We guarantee that the suffix is valid in the TypeID constructors, so this panic
	// should never be reached.
	if err != nil {
		panic(err)
	}
	return dst[:]
}

// UUID decodes the TypeID's suffix as a UUID and returns it as a hex string
func (tid TypeID) UUID() string {
	return uuid.FromBytesOrNil(tid.Bytes()).String()
}

// HasSuffix returns true if the TypeID has a non-zero suffix.
//
// This method returns false only when the suffix is the zero suffix:
// "00000000000000000000000000"
//
// Note that HasSuffix() checks only the suffix value, regardless of the prefix.
// All of these examples would return `HasSuffix() == false`:
// + "prefix_00000000000000000000000000"
// + "test_00000000000000000000000000"
// + "00000000000000000000000000"
func (tid TypeID) HasSuffix() bool {
	return tid.Suffix() != ZeroSuffix
}

// IsZero returns true if the TypeID is the zero value (empty prefix and zero suffix).
//
// Unlike HasSuffix(), IsZero() returns true only when both:
// + The prefix is empty (no type specified)
// + The suffix is the zero suffix "00000000000000000000000000"
//
// Note that the empty struct TypeID{} is encoded as the zero id.
func (tid TypeID) IsZero() bool {
	return tid.value == ""
}
