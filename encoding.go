package typeid

import (
	"encoding"
)

// TODO: Define a standardized binary encoding for typeids in the spec
// and use that to implement encoding.BinaryMarshaler and encoding.BinaryUnmarshaler

var (
	_ encoding.TextMarshaler   = (*TypeID)(nil)
	_ encoding.TextUnmarshaler = (*TypeID)(nil)
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses a TypeID from a string using the same logic as FromString()
func (tid *TypeID) UnmarshalText(text []byte) error {
	parsed, err := Parse(string(text))
	if err != nil {
		return err
	}
	*tid = parsed
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// It encodes a TypeID as a string using the same logic as String()
func (tid TypeID) MarshalText() (text []byte, err error) {
	return tid.AppendText(nil)
}

// AppendText appends the text representation of the TypeID to dst and returns
// the extended buffer.
func (tid TypeID) AppendText(dst []byte) ([]byte, error) {
	if tid.value == "" {
		return append(dst, ZeroSuffix...), nil
	}
	return append(dst, tid.value...), nil
}
