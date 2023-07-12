package typed

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses a TypeID from a string using the same logic as FromString()
func (tid *TypeID[T]) UnmarshalText(text []byte) error {
	parsed, err := FromString[T](string(text))
	if err != nil {
		return err
	}
	*tid = parsed
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// It encodes a TypeID as a string using the same logic as String()
func (tid TypeID[T]) MarshalText() (text []byte, err error) {
	encoded := tid.String()
	return []byte(encoded), nil
}
