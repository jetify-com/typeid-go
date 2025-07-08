package typeid_test

import (
	_ "embed"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/typeid/v2"
	"gopkg.in/yaml.v2"
)

// TestHasSuffix tests the HasSuffix method for various TypeID values
func TestHasSuffix(t *testing.T) {
	testdata := []struct {
		input  string
		output bool
	}{
		// HasSuffix == false values (zero suffix)
		{"00000000000000000000000000", false},
		{"prefix_00000000000000000000000000", false},
		{"other_00000000000000000000000000", false},
		// HasSuffix == true values (non-zero suffix)
		{"00000000000000000000000001", true},
		{"prefix_00000000000000000000000001", true},
		{"other_00000000000000000000000001", true},
	}

	for _, td := range testdata {
		t.Run(td.input, func(t *testing.T) {
			tid, err := typeid.Parse(td.input)
			assert.NoError(t, err)
			assert.Equal(t, td.output, tid.HasSuffix(), "TypeId.HasSuffix should be %v for id %s", td.output, td.input)
		})
	}
}

func TestRoundTrip(t *testing.T) {
	// Generate a bunch of random typeids, encode and decode from a string
	// and make sure the result is the same as the original.
	for i := 0; i < 1000; i++ {
		tid := typeid.MustGenerate("prefix")
		decoded, err := typeid.Parse(tid.String())
		assert.NoError(t, err)
		assert.Equal(t, tid, decoded)
	}

	// Repeat with the empty prefix:
	for i := 0; i < 1000; i++ {
		tid := typeid.MustGenerate("")
		decoded, err := typeid.Parse(tid.String())
		assert.NoError(t, err)
		assert.Equal(t, tid, decoded)
	}
}

//go:embed testdata/valid.yml
var validYML []byte

type ValidExample struct {
	Name   string `yaml:"name"`
	Tid    string `yaml:"typeid"`
	Prefix string `yaml:"prefix"`
	UUID   string `yaml:"uuid"`
}

func TestValidTestdata(t *testing.T) {
	assert.Greater(t, len(validYML), 0)
	var testdata []ValidExample
	err := yaml.Unmarshal(validYML, &testdata)
	assert.NoError(t, err, "Failed to unmarshal testdata")
	assert.Greater(t, len(testdata), 0)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			testValidExample(t, td)
		})
	}
}

// testValidExample tests all applicable constructors for a valid TypeID example
func testValidExample(t *testing.T, example ValidExample) {
	t.Helper()

	// Test Parse constructor
	tidParsed, err := typeid.Parse(example.Tid)
	require.NoError(t, err, "Parse should succeed for valid typeid: %s", example.Tid)

	// Test FromUUID constructor
	tidFromUUID, err := typeid.FromUUID(example.Prefix, example.UUID)
	require.NoError(t, err, "FromUUID should succeed")

	// Test FromBytes constructor
	uuidParsed, err := uuid.FromString(example.UUID)
	require.NoError(t, err, "UUID parsing should succeed")
	tidFromBytes, err := typeid.FromBytes(example.Prefix, uuidParsed.Bytes())
	require.NoError(t, err, "FromBytes should succeed")

	// All constructors should produce structurally identical TypeID objects
	assert.Equal(t, tidParsed, tidFromUUID, "Parse and FromUUID should return identical TypeID structs")
	assert.Equal(t, tidParsed, tidFromBytes, "Parse and FromBytes should return identical TypeID structs")
	assert.Equal(t, tidFromUUID, tidFromBytes, "FromUUID and FromBytes should return identical TypeID structs")

	// All constructors should produce identical string representations and components
	assert.Equal(t, example.Tid, tidParsed.String())
	assert.Equal(t, example.Tid, tidFromUUID.String())
	assert.Equal(t, example.Tid, tidFromBytes.String())

	assert.Equal(t, example.UUID, tidParsed.UUID())
	assert.Equal(t, example.UUID, tidFromUUID.UUID())
	assert.Equal(t, example.UUID, tidFromBytes.UUID())

	assert.Equal(t, example.Prefix, tidParsed.Prefix())
	assert.Equal(t, example.Prefix, tidFromUUID.Prefix())
	assert.Equal(t, example.Prefix, tidFromBytes.Prefix())

	// All constructors should produce identical suffixes
	assert.Equal(t, tidParsed.Suffix(), tidFromUUID.Suffix())
	assert.Equal(t, tidParsed.Suffix(), tidFromBytes.Suffix())
	assert.Equal(t, tidFromUUID.Suffix(), tidFromBytes.Suffix())

	// All constructors should produce identical byte arrays
	assert.Equal(t, tidParsed.Bytes(), tidFromUUID.Bytes())
	assert.Equal(t, tidParsed.Bytes(), tidFromBytes.Bytes())
	assert.Equal(t, tidFromUUID.Bytes(), tidFromBytes.Bytes())

	// All constructors should have consistent HasSuffix() behavior
	assert.Equal(t, tidParsed.HasSuffix(), tidFromUUID.HasSuffix())
	assert.Equal(t, tidParsed.HasSuffix(), tidFromBytes.HasSuffix())
	assert.Equal(t, tidFromUUID.HasSuffix(), tidFromBytes.HasSuffix())
}

// invalidPrefixTestCases contains all invalid prefix cases that should be rejected
// by all TypeID constructors. These cases are based on the invalid.yml fixture.
var invalidPrefixTestCases = []struct {
	name   string
	prefix string
	desc   string
}{
	{
		name:   "prefix-uppercase",
		prefix: "PREFIX",
		desc:   "prefix with uppercase letters",
	},
	{
		name:   "prefix-numeric",
		prefix: "12345",
		desc:   "prefix with only numbers",
	},
	{
		name:   "prefix-period",
		prefix: "pre.fix",
		desc:   "prefix with period/dot",
	},
	{
		name:   "prefix-non-ascii",
		prefix: "préfix",
		desc:   "prefix with non-ASCII characters",
	},
	{
		name:   "prefix-spaces",
		prefix: "  prefix",
		desc:   "prefix with leading spaces",
	},
	{
		name:   "prefix-space-middle",
		prefix: "pre fix",
		desc:   "prefix with space in middle",
	},
	{
		name:   "prefix-64-chars",
		prefix: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl",
		desc:   "prefix with 64 characters (too long)",
	},
	{
		name:   "prefix-underscore-start",
		prefix: "_prefix",
		desc:   "prefix starting with underscore",
	},
	{
		name:   "prefix-underscore-end",
		prefix: "prefix_",
		desc:   "prefix ending with underscore",
	},
	{
		name:   "prefix-hyphen",
		prefix: "pre-fix",
		desc:   "prefix with hyphen",
	},
	{
		name:   "prefix-special-chars",
		prefix: "pre@fix",
		desc:   "prefix with special characters",
	},
}

func TestInvalidPrefix(t *testing.T) {
	for _, tc := range invalidPrefixTestCases {
		t.Run(tc.name, func(t *testing.T) {
			testInvalidPrefix(t, tc.prefix, tc.desc)
		})
	}
}

//go:embed testdata/invalid.yml
var invalidYML []byte

type InvalidExample struct {
	Name string `yaml:"name"`
	Tid  string `yaml:"typeid"`
}

func TestInvalidTestdata(t *testing.T) {
	assert.Greater(t, len(invalidYML), 0)
	var testdata []InvalidExample
	err := yaml.Unmarshal(invalidYML, &testdata)
	assert.NoError(t, err, "Failed to unmarshal testdata")
	assert.Greater(t, len(testdata), 0)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			testInvalidExample(t, td)
		})
	}
}

// testInvalidExample tests that all constructors properly reject invalid TypeIDs
func testInvalidExample(t *testing.T, example InvalidExample) {
	t.Helper()

	// Test Parse constructor - should always fail for invalid examples
	_, err := typeid.Parse(example.Tid)
	assert.Error(t, err, "Parse should fail for invalid typeid: %s", example.Tid)
}

// testInvalidPrefix tests that all constructors properly reject invalid prefixes
func testInvalidPrefix(t *testing.T, prefix string, desc string) {
	t.Helper()

	// Test Generate constructor with invalid prefix
	_, err := typeid.Generate(prefix)
	assert.Error(t, err, "Generate should fail with %s", desc)

	// Test MustGenerate constructor with invalid prefix - should panic
	assert.Panics(t, func() {
		typeid.MustGenerate(prefix)
	}, "MustGenerate should panic with %s", desc)

	// Test FromUUID constructor with invalid prefix
	// Use a valid UUID to isolate prefix validation
	validUUID := "00000000-0000-0000-0000-000000000000"
	_, err = typeid.FromUUID(prefix, validUUID)
	assert.Error(t, err, "FromUUID should fail with %s", desc)

	// Test FromBytes constructor with invalid prefix
	zeroBytes := make([]byte, 16)
	_, err = typeid.FromBytes(prefix, zeroBytes)
	assert.Error(t, err, "FromBytes should fail with %s", desc)

	// Test Parse with a complete TypeID string containing the invalid prefix
	if prefix != "" {
		invalidTypeID := prefix + "_00000000000000000000000000"
		_, err = typeid.Parse(invalidTypeID)
		assert.Error(t, err, "Parse should fail with %s", desc)
	}
}

// TestZero verifies that all constructors handle zero TypeIDs consistently.
// This test ensures that Parse, FromUUID, and FromBytes all return identical zero TypeIDs
// that are equal to the zero-value struct and the canonical ZeroID.
func TestZero(t *testing.T) {
	// The nil UUID string representation
	nilUUID := "00000000-0000-0000-0000-000000000000"
	nilTypeID := "00000000000000000000000000"

	// Zero-value struct (uninitialized)
	var zeroValue typeid.TypeID

	// Parse zero TypeID
	tidParsed, err := typeid.Parse(nilTypeID)
	require.NoError(t, err)

	// FromUUID with zero UUID
	tidFromUUID, err := typeid.FromUUID("", nilUUID)
	require.NoError(t, err)

	// FromBytes with zero bytes
	zeroBytes := make([]byte, 16)
	tidFromBytes, err := typeid.FromBytes("", zeroBytes)
	require.NoError(t, err)

	// All should be structurally identical
	assert.Equal(t, zeroValue, tidParsed, "Zero-value struct should equal parsed zero TypeID")
	assert.Equal(t, zeroValue, tidFromUUID, "Zero-value struct should equal FromUUID zero TypeID")
	assert.Equal(t, zeroValue, tidFromBytes, "Zero-value struct should equal FromBytes zero TypeID")
	assert.Equal(t, tidParsed, tidFromUUID, "Parse and FromUUID should return identical zero TypeIDs")
	assert.Equal(t, tidParsed, tidFromBytes, "Parse and FromBytes should return identical zero TypeIDs")
	assert.Equal(t, tidFromUUID, tidFromBytes, "FromUUID and FromBytes should return identical zero TypeIDs")

	// All should have identical functional behavior
	assert.Equal(t, nilTypeID, zeroValue.String())
	assert.Equal(t, nilTypeID, tidParsed.String())
	assert.Equal(t, nilTypeID, tidFromUUID.String())
	assert.Equal(t, nilTypeID, tidFromBytes.String())

	assert.False(t, zeroValue.HasSuffix())
	assert.False(t, tidParsed.HasSuffix())
	assert.False(t, tidFromUUID.HasSuffix())
	assert.False(t, tidFromBytes.HasSuffix())

	// All should return true for IsZero() since they have empty prefix and zero suffix
	assert.True(t, zeroValue.IsZero())
	assert.True(t, tidParsed.IsZero())
	assert.True(t, tidFromUUID.IsZero())
	assert.True(t, tidFromBytes.IsZero())

	assert.Equal(t, nilUUID, zeroValue.UUID())
	assert.Equal(t, nilUUID, tidParsed.UUID())
	assert.Equal(t, nilUUID, tidFromUUID.UUID())
	assert.Equal(t, nilUUID, tidFromBytes.UUID())

	// All should have identical Bytes() output
	assert.Equal(t, zeroValue.Bytes(), tidParsed.Bytes())
	assert.Equal(t, zeroValue.Bytes(), tidFromUUID.Bytes())
	assert.Equal(t, zeroValue.Bytes(), tidFromBytes.Bytes())
}

// TestIsZero tests the IsZero method for various TypeID values
func TestIsZero(t *testing.T) {
	testdata := []struct {
		input  string
		output bool
	}{
		// IsZero == true values (empty prefix AND zero suffix)
		{"00000000000000000000000000", true},
		// IsZero == false values (has prefix OR non-zero suffix)
		{"prefix_00000000000000000000000000", false}, // has prefix
		{"other_00000000000000000000000000", false},  // has prefix
		{"00000000000000000000000001", false},        // no prefix but non-zero suffix
		{"prefix_00000000000000000000000001", false}, // has prefix and non-zero suffix
		{"other_00000000000000000000000001", false},  // has prefix and non-zero suffix
	}

	for _, td := range testdata {
		t.Run(td.input, func(t *testing.T) {
			tid, err := typeid.Parse(td.input)
			assert.NoError(t, err)
			assert.Equal(t, td.output, tid.IsZero(), "TypeId.IsZero should be %v for id %s", td.output, td.input)
		})
	}
}
