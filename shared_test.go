package typeid

// MustParse returns a TypeID if the error is nil, otherwise panics.
// Used in tests to create a TypeID in a single line as follows:
// tid := MustParse("prefix_abc123")
func MustParse(s string) TypeID {
	tid, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return tid
}
