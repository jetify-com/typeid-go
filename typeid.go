package typeid

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"go.jetpack.io/typeid/base32"
)

type TypeID struct {
	prefix string
	suffix string
}

var Nil = TypeID{
	prefix: "",
	suffix: "00000000000000000000000000",
}

func New(prefix string) (TypeID, error) {
	return From(prefix, "")
}

func (tid TypeID) Type() string {
	return tid.prefix
}

func (tid TypeID) Suffix() string {
	return tid.suffix
}

func (tid TypeID) String() string {
	if tid.prefix == "" {
		return tid.suffix
	}
	return tid.prefix + "_" + tid.Suffix()
}

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

func (tid TypeID) UUID() string {
	return uuid.FromBytesOrNil(tid.UUIDBytes()).String()
}

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

func FromString(s string) (TypeID, error) {
	switch parts := strings.SplitN(s, "_", 2); len(parts) {
	case 1:
		return From("", parts[0])
	case 2:
		return From(parts[0], parts[1])
	default:
		return Nil, fmt.Errorf("invalid typeid: %s", s)
	}
}

func FromUUID(prefix string, uidStr string) (TypeID, error) {
	uid, err := uuid.FromString(uidStr)
	if err != nil {
		return Nil, err
	}
	suffix := base32.Encode(uid)
	return From(prefix, suffix)
}

func FromUUIDBytes(prefix string, bytes []byte) (TypeID, error) {
	uidStr := uuid.FromBytesOrNil(bytes).String()
	return FromUUID(prefix, uidStr)
}

func Must(tid TypeID, err error) TypeID {
	if err != nil {
		panic(err)
	}
	return tid
}

func validatePrefix(prefix string) error {
	// Ensure that the prefix only has lowercase ASCII characters
	for _, c := range prefix {
		if c < 'a' || c > 'z' {
			return fmt.Errorf("invalid prefix: '%s'. Prefix should match [a-z]+", prefix)
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
