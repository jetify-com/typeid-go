package typeid_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.jetpack.io/typeid"
)

func TestScan(t *testing.T) {
	testdata := []struct {
		name     string
		input    any
		expected typeid.AnyID
	}{
		{"valid", "prefix_00041061050r3gg28a1c60t3gf", typeid.Must(typeid.FromString("prefix_00041061050r3gg28a1c60t3gf"))},
		{"nil", nil, typeid.AnyID{}},
		{"empty string", "", typeid.AnyID{}},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			var scanned typeid.AnyID
			err := scanned.Scan(td.input)
			assert.NoError(t, err)

			assert.Equal(t, td.expected, scanned)
			assert.Equal(t, td.expected.String(), scanned.String())
		})
	}
}

func TestValuer(t *testing.T) {
	expected := "prefix_00041061050r3gg28a1c60t3gf"
	tid := typeid.Must(typeid.FromString("prefix_00041061050r3gg28a1c60t3gf"))
	actual, err := tid.Value()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
