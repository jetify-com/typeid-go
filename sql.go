package typeid

import (
	"database/sql/driver"
	"fmt"
)

// TODO: decide if we want nullable (or just use pointers)

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

// NullableID is wrapper for nullable columns.
type NullableID struct {
	TypeID TypeID
	Valid  bool
}

func (n NullableID) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil // SQL NULL
	}
	return n.TypeID.Value() // Delegate to TypeID
}

func (n *NullableID) Scan(src any) error {
	if src == nil {
		n.TypeID, n.Valid = zeroID, false
		return nil
	}

	// Empty string is invalid even for nullable columns - force explicit NULL usage
	if str, ok := src.(string); ok && str == "" {
		return &validationError{
			Message: "empty string is invalid TypeID",
		}
	}

	// Try to scan the TypeID, only set Valid=true if successful
	err := n.TypeID.Scan(src)
	if err != nil {
		n.Valid = false
		return err
	}

	n.Valid = true
	return nil
}
