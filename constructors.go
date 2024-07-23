package typeid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
	"go.jetify.com/typeid/base32"
)

var ErrConstructor = errors.New("constructor error")

// New returns a new TypeID of the given type with a random suffix.
//
// Use the generic argument to pass in your typeid Subtype:
//
// Example:
//
//	  type UserID struct {
//		   typeid.TypeID[UserPrefix]
//	  }
//	  id, err := typeid.New[UserID]()
func New[T Subtype, PT SubtypePtr[T]]() (T, error) {
	if isAnyID[T]() {
		var id T
		return id, errors.New("constructor error: use WithPrefix(), New() is for Subtypes")
	}

	prefix := defaultType[T]()
	return from[T, PT](prefix, "")
}

// WithPrefix returns a new TypeID with the given prefix and a random suffix.
// If you want to create an id without a prefix, pass an empty string.
func WithPrefix(prefix string) (AnyID, error) {
	return from[AnyID](prefix, "")
}

// From returns a new TypeID with the given prefix and suffix.
// If suffix is the empty string, a random suffix will be generated.
// If you want to create an id without a prefix, pass an empty string as the prefix.
func From(prefix string, suffix string) (AnyID, error) {
	return from[AnyID](prefix, suffix)
}

// FromSuffix returns a new TypeID of the given suffix and type. The prefix
// is inferred from the Subtype.
//
// Example:
//
//	  type UserID struct {
//		   typeid.TypeID[UserPrefix]
//	  }
//	  id, err := typeid.FromSuffix[UserID]("00041061050r3gg28a1c60t3gf")
func FromSuffix[T Subtype, PT SubtypePtr[T]](suffix string) (T, error) {
	if isAnyID[T]() {
		var id T
		return id, errors.New("constructor error: use From(prefix, suffix), FromSuffix is for Subtypes")
	}

	prefix := defaultType[T]()
	return parse[T, PT](prefix, suffix)
}

// FromString parses a TypeID from a string of the form <prefix>_<suffix>
func FromString(s string) (AnyID, error) {
	return Parse[AnyID](s)
}

// Parse parses a TypeID from a string of the form <prefix>_<suffix>
// and ensures the TypeID is of the right type.
//
// Example:
//
//	  type UserID struct {
//		   typeid.TypeID[UserPrefix]
//	  }
//	  id, err := typeid.Parse[UserID]("user_00041061050r3gg28a1c60t3gf")
func Parse[T Subtype, PT SubtypePtr[T]](s string) (T, error) {
	prefix, suffix, err := split(s)
	if err != nil {
		var id T
		return id, err
	}
	return parse[T, PT](prefix, suffix)
}

func split(id string) (string, string, error) {
	index := strings.LastIndex(id, "_")
	if index == -1 {
		return "", id, nil
	}

	prefix := id[:index]
	suffix := id[index+1:]
	if prefix == "" {
		return "", "", errors.New("prefix cannot be empty when there's a separator")
	}
	return prefix, suffix, nil
}

// FromUUID encodes the given UUID (in hex string form) as a TypeID
func FromUUID[T Subtype, PT SubtypePtr[T]](uidStr string) (T, error) {
	if isAnyID[T]() {
		var id T
		return id, fmt.Errorf(
			"%w: use FromUUIDWithPrefix(), FromUUID() is for Subtypes",
			ErrConstructor,
		)
	}
	return fromUUID[T, PT](defaultPrefix[T](), uidStr)
}

// FromUUIDBytes encodes the given UUID (in byte form) as a TypeID
func FromUUIDBytes[T Subtype, PT SubtypePtr[T]](bytes []byte) (T, error) {
	if isAnyID[T]() {
		var id T
		return id, fmt.Errorf(
			"%w: use FromUUIDBytesWithPrefix(), FromUUIDBytes() is for Subtypes",
			ErrConstructor,
		)
	}
	uidStr := uuid.FromBytesOrNil(bytes).String()
	return FromUUID[T, PT](uidStr)
}

// FromUUIDWithPrefix encodes the given UUID (in hex string form) as a TypeID
// with the given prefix.
func FromUUIDWithPrefix(prefix string, uidStr string) (AnyID, error) {
	return fromUUID[AnyID](prefix, uidStr)
}

// FromUUID encodes the given UUID (in byte form) as a TypeID with the given
// prefix.
func FromUUIDBytesWithPrefix(prefix string, bytes []byte) (AnyID, error) {
	uidStr := uuid.FromBytesOrNil(bytes).String()
	return FromUUIDWithPrefix(prefix, uidStr)
}

func fromUUID[T Subtype, PT SubtypePtr[T]](prefix, uidStr string) (T, error) {
	uid, err := uuid.FromString(uidStr)
	var nilID T

	if err != nil {
		return nilID, err
	}
	suffix := base32.Encode(uid)
	return parse[T, PT](prefix, suffix)
}

func parse[T Subtype, PT SubtypePtr[T]](prefix string, suffix string) (T, error) {
	if suffix == "" {
		var id T
		return id, errors.New("suffix can't be the empty string")
	}
	return from[T, PT](prefix, suffix)
}

func from[T Subtype, PT SubtypePtr[T]](prefix string, suffix string) (T, error) {
	var tid T
	if err := validatePrefix[T](prefix); err != nil {
		return tid, err
	}

	if suffix == "" {
		uid, err := uuid.NewV7()
		if err != nil {
			return tid, err
		}
		suffix = base32.Encode(uid)
	}

	if err := validateSuffix(suffix); err != nil {
		return tid, err
	}

	PT(&tid).init(prefix, suffix)
	return tid, nil
}
