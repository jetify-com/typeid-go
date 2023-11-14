package typeid

import (
	"database/sql/driver"
	"fmt"
)

// Scan implements the sql.Scanner interface so the TypeIDs can be read from
// databases transparently. Currently database types that map to string are
// supported.
func (tid *TypeID[P]) Scan(src any) error {
	switch obj := src.(type) {
	case nil:
		return nil
	case string:
		if src == "" {
			return nil
		}
		return tid.UnmarshalText([]byte(obj))
	// TODO: add supporte for []byte
	// we don't just want to store the full string as a byte array. Instead
	// we should encode using the UUID bytes. We could add support for
	// Binary Marshalling and Unmarshalling at the same time.
	default:
		return fmt.Errorf("unsupported scan type %T", obj)
	}
}

// Value implements the sql.Valuer interface so that TypeIDs can be written
// to databases transparently. Currently, TypeIDs map to strings.
func (tid TypeID[P]) Value() (driver.Value, error) {
	return tid.String(), nil
}
