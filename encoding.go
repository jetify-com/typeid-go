package typeid

import "encoding"

// TODO: Define a standardized binary encoding for typeids in the spec
// and use that to implement encoding.BinaryMarshaler and encoding.BinaryUnmarshaler

var _ encoding.TextMarshaler = (*TypeID)(nil)
var _ encoding.TextUnmarshaler = (*TypeID)(nil)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses a TypeID from a string using the same logic as FromString()
func (tid *TypeID) UnmarshalText(text []byte) error {
	parsed, err := FromString(string(text))
	if err != nil {
		return err
	}
	*tid = parsed
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// It encodes a TypeID as a string using the same logic as String()
func (tid TypeID) MarshalText() (text []byte, err error) {
	encoded := tid.String()
	return []byte(encoded), nil
}
