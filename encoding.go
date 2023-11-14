package typeid

import (
	"encoding"
)

// TODO: Define a standardized binary encoding for typeids in the spec
// and use that to implement encoding.BinaryMarshaler and encoding.BinaryUnmarshaler

var _ encoding.TextMarshaler = (*TypeID[AnyPrefix])(nil)
var _ encoding.TextUnmarshaler = (*TypeID[AnyPrefix])(nil)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses a TypeID from a string using the same logic as FromString()
func (tid *TypeID[P]) UnmarshalText(text []byte) error {
	parsed, err := Parse[TypeID[P]](string(text))
	if err != nil {
		return err
	}
	*tid = parsed
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// It encodes a TypeID as a string using the same logic as String()
func (tid TypeID[P]) MarshalText() (text []byte, err error) {
	encoded := tid.String()
	return []byte(encoded), nil
}
