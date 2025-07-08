package typeid

import (
	"database/sql/driver"
	"fmt"
)

// For nullable TypeID columns, use sql.Null[TypeID].

// Scan implements the sql.Scanner interface so the TypeIDs can be read from
// databases transparently. Currently database types that map to string are
// supported.
func (tid *TypeID) Scan(src any) error {
	switch obj := src.(type) {
	case nil:
		return &validationError{
			Message: "cannot scan NULL into TypeID",
		}
	case string:
		if obj == "" {
			return &validationError{
				Message: "cannot scan empty string into TypeID",
			}
		}
		return tid.UnmarshalText([]byte(obj))
	default:
		return &validationError{
			Message: fmt.Sprintf("unsupported scan type %T", obj),
		}
	}
}

// Value implements the sql.Valuer interface so that TypeIDs can be written
// to databases transparently. Currently, TypeIDs map to strings.
func (tid TypeID) Value() (driver.Value, error) {
	return tid.String(), nil
}
