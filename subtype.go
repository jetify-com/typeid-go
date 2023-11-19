package typeid

// Subtype is an interface used to create a more specific subtype of TypeID
// For example, if you want to create an `OrgID` type that only accepts
// an `org_` prefix.
type Subtype interface {
	Prefix() string
	Suffix() string
	String() string
	UUIDBytes() []byte
	UUID() string

	isTypeID() bool
}

var _ Subtype = (*TypeID[AnyPrefix])(nil)

type SubtypePtr[T any] interface {
	*T
	init(prefix string, suffix string)
}

func (tid *TypeID[P]) init(prefix string, suffix string) {
	// In general TypeID is an immutable value-type, and pretty much every
	// "mutation" should return a copy with the modifications instead of modifying
	// the original. We make an exception for this *private* method, because
	// sometimes we need to modify the fields in the process of initializing
	// a new subtype.

	// Only store the prefix if dealing with a subtype:
	if isAnyPrefix[P]() {
		tid.prefix = prefix
	}

	// If we're dealing with the "nil" suffix, we don't need to store it.
	if suffix != nilSuffix {
		tid.suffix = suffix
	}
}

func (tid TypeID[P]) isTypeID() bool {
	return true
}

func isAnyID[T Subtype]() bool {
	var id T
	switch any(id).(type) {
	case TypeID[AnyPrefix]:
		return true
	case AnyID:
		return true
	default:
		return false
	}
}

func defaultType[T Subtype]() string {
	var id T
	return id.Prefix()
}
