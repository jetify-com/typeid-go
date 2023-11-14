package typeid

// PrefixType is the interface that defines the type if a type id.
// Implement your own version of this interface if you want to define a custom
// type:
// type UserPrefix struct {}
// func (UserPrefix) Prefix() string { return "user" }
type PrefixType interface {
	Prefix() string
}

// Any is a special prefix that can be used to represent TypeIDs that allow for
// any valid prefix.
type AnyPrefix struct{}

func (a AnyPrefix) Prefix() string {
	return "*" // Any is treated specially, so in practice this string will never be used.
}

// AnyID represents TypeIDs that accept any valid prefix.
type AnyID struct {
	TypeID[AnyPrefix]
}

func isAnyPrefix[P PrefixType]() bool {
	var prefixType P
	switch any(prefixType).(type) {
	case AnyPrefix:
		return true
	default:
		return false
	}
}

func defaultPrefix[P PrefixType]() string {
	var prefixType P
	return prefixType.Prefix()
}
