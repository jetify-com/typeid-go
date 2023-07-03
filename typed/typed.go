package typed

import (
	"fmt"

	untyped "go.jetpack.io/typeid"
)

// TypePrefix is an interface used to represent a statically checked type prefix.
// Example:
// type userPrefix struct{}
// func (userPrefix) Type() string { return "user" }
//
//	type UserID struct {
//		typeid.TypeID[userPrefix]
//	}
type TypePrefix interface {
	Type() string
}

// TypeID is a unique identifier with a given type as defined by the TypeID spec
type TypeID[T TypePrefix] untyped.TypeID

// New returns a new TypeID with a random suffix and the given type.
func New[T TypePrefix]() (TypeID[T], error) {
	tid, err := untyped.New(Type[T]())
	if err != nil {
		// Clients should ignore the id value when an error is present, but just
		// in case, construct a "nil" id of the given type.
		return Nil[T](), err
	}
	return (TypeID[T])(tid), nil
}

func Type[T TypePrefix]() string {
	var prefix T
	return prefix.Type()
}

// Nil returns the null typeid of the given type.
func Nil[T TypePrefix]() TypeID[T] {
	return TypeID[T](untyped.Must(untyped.From(Type[T](), "00000000000000000000000000")))
}

// Type returns the type prefix of the TypeID
func (tid TypeID[T]) Type() string {
	var prefix T
	return prefix.Type()
}

// Suffix returns the suffix of the TypeID in it's canonical base32 representation.
func (tid TypeID[T]) Suffix() string {
	return untyped.TypeID(tid).Suffix()
}

// String returns the TypeID in it's canonical string representation of the form:
// <prefix>_<suffix> where <suffix> is the canonical base32 representation of the UUID
func (tid TypeID[T]) String() string {
	return untyped.TypeID(tid).String()
}

// UUIDBytes decodes the TypeID's suffix as a UUID and returns it's bytes
func (tid TypeID[T]) UUIDBytes() []byte {
	return untyped.TypeID(tid).UUIDBytes()
}

// UUID decode the TypeID's suffix as a UUID and returns it as a hex formatted string
func (tid TypeID[T]) UUID() string {
	return untyped.TypeID(tid).UUID()
}

// From returns a new TypeID of the given type using the provided suffix
func From[T TypePrefix](suffix string) (TypeID[T], error) {
	tid, err := untyped.From(Type[T](), suffix)
	if err != nil {
		return Nil[T](), err
	}
	return (TypeID[T])(tid), nil
}

// FromString parses a TypeID from the given string. Returns an error if the
// string is not a valid TypeID, OR if the type prefix does not match the
// expected type.
func FromString[T TypePrefix](s string) (TypeID[T], error) {
	tid, err := untyped.FromString(s)
	if err != nil {
		return Nil[T](), err
	}
	if tid.Type() != Type[T]() {
		return Nil[T](), fmt.Errorf("invalid type, expected %s but got %s", Type[T](), tid.Type())
	}
	return (TypeID[T])(tid), nil
}

// FromUUID returns a new TypeID of the given type using the provided UUID
func FromUUID[T TypePrefix](uuid string) (TypeID[T], error) {
	tid, err := untyped.FromUUID(Type[T](), uuid)
	if err != nil {
		return Nil[T](), err
	}
	return (TypeID[T])(tid), nil
}

// FromUUIDBytes returns a new TypeID of the given type using the provided UUID bytes
func FromUUIDBytes[T TypePrefix](uuid []byte) (TypeID[T], error) {
	tid, err := untyped.FromUUIDBytes(Type[T](), uuid)
	if err != nil {
		return Nil[T](), err
	}
	return (TypeID[T])(tid), nil
}

// Must panics if the given error is non-nil, otherwise it returns the given TypeID
func Must[T TypePrefix](tid TypeID[T], err error) TypeID[T] {
	if err != nil {
		panic(err)
	}
	return tid
}
