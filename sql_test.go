package typeid_test

import (
	"database/sql"
	_ "embed"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/typeid"
	"gopkg.in/yaml.v2"
)

//go:embed testdata/valid.yml
var validSQLYML []byte

//go:embed testdata/invalid.yml
var invalidSQLYML []byte

func TestScanValid(t *testing.T) {
	var testdata []ValidExample
	err := yaml.Unmarshal(validSQLYML, &testdata)
	require.NoError(t, err)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			// Test Scan with string input
			var scanned typeid.TypeID
			err := scanned.Scan(td.Tid)
			assert.NoError(t, err)

			expected := typeid.MustParse(td.Tid)
			assert.Equal(t, expected, scanned)
			assert.Equal(t, td.Tid, scanned.String())
		})
	}
}

func TestScanSpecialCases(t *testing.T) {
	testdata := []struct {
		name        string
		input       any
		expectError bool
	}{
		{"nil", nil, true},
		{"empty string", "", true},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			var scanned typeid.TypeID
			err := scanned.Scan(td.input)
			if td.expectError {
				assert.Error(t, err)
				// Verify that scan errors are validation errors
				assert.True(t, errors.Is(err, typeid.ErrValidation), "expected validation error")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestScanInvalid(t *testing.T) {
	var testdata []InvalidExample
	err := yaml.Unmarshal(invalidSQLYML, &testdata)
	require.NoError(t, err)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			// Test Scan with invalid TypeID strings
			var scanned typeid.TypeID
			err := scanned.Scan(td.Tid)
			assert.Error(t, err, "Scan should fail for invalid typeid: %s", td.Tid)
		})
	}
}

func TestScanUnsupportedType(t *testing.T) {
	testdata := []struct {
		name  string
		input any
	}{
		{"int", 123},
		{"float64", 123.45},
		{"bool", true},
		{"[]byte", []byte("test")},
		{"time.Time", time.Now()},
		{"struct", struct{ field string }{field: "test"}},
		{"map", map[string]string{"key": "value"}},
		{"slice", []string{"a", "b", "c"}},
		{"int64", int64(123)},
		{"uint", uint(123)},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			var scanned typeid.TypeID
			err := scanned.Scan(td.input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "unsupported scan type")
			// Verify that scan errors are validation errors
			assert.True(t, errors.Is(err, typeid.ErrValidation), "expected validation error")
		})
	}
}

func TestValue(t *testing.T) {
	var testdata []ValidExample
	err := yaml.Unmarshal(validSQLYML, &testdata)
	require.NoError(t, err)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			tid := typeid.MustParse(td.Tid)
			actual, err := tid.Value()
			assert.NoError(t, err)
			assert.Equal(t, td.Tid, actual)
		})
	}
}

// Test sql.Null[TypeID] to verify it works for nullable database columns
func TestSQLNullScanValid(t *testing.T) {
	var testdata []ValidExample
	err := yaml.Unmarshal(validSQLYML, &testdata)
	require.NoError(t, err)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			// Test sql.Null[TypeID].Scan with valid TypeID strings
			var scanned sql.Null[typeid.TypeID]
			err := scanned.Scan(td.Tid)
			assert.NoError(t, err)

			expected := typeid.MustParse(td.Tid)
			assert.True(t, scanned.Valid, "sql.Null[TypeID] should be valid for valid typeid")
			assert.Equal(t, expected, scanned.V)
			assert.Equal(t, td.Tid, scanned.V.String())
		})
	}
}

func TestSQLNullScanSpecialCases(t *testing.T) {
	testdata := []struct {
		name        string
		input       any
		expected    sql.Null[typeid.TypeID]
		expectError bool
	}{
		{"nil", nil, sql.Null[typeid.TypeID]{Valid: false}, false},
		{"empty string", "", sql.Null[typeid.TypeID]{}, true},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			var scanned sql.Null[typeid.TypeID]
			err := scanned.Scan(td.input)

			if td.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot scan empty string into TypeID")
				// Verify that scan errors are validation errors
				assert.True(t, errors.Is(err, typeid.ErrValidation), "expected validation error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, td.expected.Valid, scanned.Valid)
				if td.expected.Valid {
					assert.Equal(t, td.expected.V, scanned.V)
				}
			}
		})
	}
}

func TestSQLNullValue(t *testing.T) {
	// Test the invalid case (Valid: false)
	t.Run("invalid", func(t *testing.T) {
		invalid := sql.Null[typeid.TypeID]{Valid: false}
		actual, err := invalid.Value()
		assert.NoError(t, err)
		assert.Equal(t, nil, actual)
	})

	// Test all valid examples from YAML
	var testdata []ValidExample
	err := yaml.Unmarshal(validSQLYML, &testdata)
	require.NoError(t, err)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			tid := typeid.MustParse(td.Tid)
			nullable := sql.Null[typeid.TypeID]{V: tid, Valid: true}
			actual, err := nullable.Value()
			assert.NoError(t, err)
			assert.Equal(t, td.Tid, actual)
		})
	}
}

func TestSQLNullScanInvalid(t *testing.T) {
	var testdata []InvalidExample
	err := yaml.Unmarshal(invalidSQLYML, &testdata)
	require.NoError(t, err)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			// Test sql.Null[TypeID].Scan with invalid TypeID strings
			var scanned sql.Null[typeid.TypeID]
			err := scanned.Scan(td.Tid)
			assert.Error(t, err, "sql.Null[TypeID].Scan should fail for invalid typeid: %s", td.Tid)
		})
	}
}

func TestSQLNullScanUnsupportedType(t *testing.T) {
	testdata := []struct {
		name  string
		input any
	}{
		{"int", 123},
		{"float64", 123.45},
		{"bool", true},
		{"[]byte", []byte("test")},
		{"time.Time", time.Now()},
		{"struct", struct{ field string }{field: "test"}},
		{"map", map[string]string{"key": "value"}},
		{"slice", []string{"a", "b", "c"}},
		{"int64", int64(123)},
		{"uint", uint(123)},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			var scanned sql.Null[typeid.TypeID]
			err := scanned.Scan(td.input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "unsupported scan type")
			assert.True(t, errors.Is(err, typeid.ErrValidation), "expected validation error")
		})
	}
}
