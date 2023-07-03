package typeid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
	"go.jetpack.io/typeid/base32"
)

// TypeID is a unique identifier with a given type as defined by the TypeID spec
type TypeID struct {
	prefix string
	suffix string
}

// Nil represents an the null TypeID
var Nil = TypeID{
	prefix: "",
	suffix: "00000000000000000000000000",
}

// New returns a new TypeID with the given prefix and a random suffix.
// If you want to create an id without a prefix, pass an empty string.
func New(prefix string) (TypeID, error) {
	return From(prefix, "")
}

// Type returns the type prefix of the TypeID
func (tid TypeID) Type() string {
	return tid.prefix
}

// Suffix returns the suffix of the TypeID in it's canonical base32 representation.
func (tid TypeID) Suffix() string {
	return tid.suffix
}

// String returns the TypeID in it's canonical string representation of the form:
// <prefix>_<suffix> where <suffix> is the canonical base32 representation of the UUID
func (tid TypeID) String() string {
	if tid.prefix == "" {
		return tid.suffix
	}
	return tid.prefix + "_" + tid.Suffix()
}

// UUIDBytes decodes the TypeID's suffix as a UUID and returns it's bytes
func (tid TypeID) UUIDBytes() []byte {
	b, err := base32.Decode(tid.suffix)

	// Decode only fails if the suffix cannot be decoded for one of two reasons:
	// 1. The suffix is not 26 characters long
	// 2. The suffix contains characters that are not in the base32 alphabet
	// We gurantee that the suffix is valid in the TypeID constructors, so this panic
	// should never be reached.
	if err != nil {
		panic(err)
	}
	return b
}

// UUID decodes the TypeID's suffix as a UUID and returns it as a hex string
func (tid TypeID) UUID() string {
	return uuid.FromBytesOrNil(tid.UUIDBytes()).String()
}

// From returns a new TypeID with the given prefix and suffix.
// If suffix is the empty string, a random suffix will be generated.
// If you want to create an id without a prefix, pass an empty string as the prefix.
func From(prefix string, suffix string) (TypeID, error) {
	if err := validatePrefix(prefix); err != nil {
		return Nil, err
	}

	if suffix == "" {
		uid, err := uuid.NewV7()
		if err != nil {
			return Nil, err
		}
		suffix = base32.Encode(uid)
	}

	if err := validateSuffix(suffix); err != nil {
		return Nil, err
	}

	return TypeID{
		prefix: prefix,
		suffix: suffix,
	}, nil

}

// FromString parses a TypeID from a string of the form <prefix>_<suffix>
func FromString(s string) (TypeID, error) {
	switch parts := strings.SplitN(s, "_", 2); len(parts) {
	case 1:
		return From("", parts[0])
	case 2:
		if parts[0] == "" {
			return Nil, errors.New("prefix cannot be empty when there's a separator")
		}
		return From(parts[0], parts[1])
	default:
		return Nil, fmt.Errorf("invalid typeid: %s", s)
	}
}

// FromUUID encodes the given UUID (in hex string form) as a TypeID with the given prefix.
func FromUUID(prefix string, uidStr string) (TypeID, error) {
	uid, err := uuid.FromString(uidStr)
	if err != nil {
		return Nil, err
	}
	suffix := base32.Encode(uid)
	return From(prefix, suffix)
}

// FromUUID encodes the given UUID (in byte form) as a TypeID with the given prefix.
func FromUUIDBytes(prefix string, bytes []byte) (TypeID, error) {
	uidStr := uuid.FromBytesOrNil(bytes).String()
	return FromUUID(prefix, uidStr)
}

// Must returns a TypeID if the error is nil, otherwise panics.
// Often used with New() to create a TypeID in a single line as follows:
// tid := Must(New("prefix"))
func Must(tid TypeID, err error) TypeID {
	if err != nil {
		panic(err)
	}
	return tid
}

func validatePrefix(prefix string) error {
	if prefix == "" {
		return nil
	}

	if len(prefix) > 63 {
		return fmt.Errorf("invalid prefix: %s. Prefix length is %d, expected <= 63", prefix, len(prefix))
	}

	// Ensure that the prefix only has lowercase ASCII characters
	for _, c := range prefix {
		if c < 'a' || c > 'z' {
			return fmt.Errorf("invalid prefix: '%s'. Prefix should match [a-z]{0,63}", prefix)
		}
	}
	return nil
}

func validateSuffix(suffix string) error {
	// Validate the suffix by decoding it:
	// 1. If the suffix is empty, it is valid
	// 2. If the suffix is not empty, it must be a valid base32 string
	if suffix == "" {
		return nil
	}
	if _, err := base32.Decode(suffix); err != nil {
		return fmt.Errorf("invalid suffix: %w", err)
	}
	return nil
}
