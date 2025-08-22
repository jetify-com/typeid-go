package typeid_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/typeid/v2"
)

//go:embed testdata/valid.yml
var validEncodingYML []byte

//go:embed testdata/invalid.yml
var invalidEncodingYML []byte

func TestJSONValid(t *testing.T) {
	var testdata []ValidExample
	err := yaml.Unmarshal(validEncodingYML, &testdata)
	require.NoError(t, err)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			// Test MarshalText via JSON encoding
			tid := typeid.MustParse(td.Tid)
			encoded, err := json.Marshal(tid)
			assert.NoError(t, err)
			assert.Equal(t, `"`+td.Tid+`"`, string(encoded))

			// Test UnmarshalText via JSON decoding
			var decoded typeid.TypeID
			err = json.Unmarshal(encoded, &decoded)
			assert.NoError(t, err)
			assert.Equal(t, tid, decoded)
			assert.Equal(t, td.Tid, decoded.String())
		})
	}
}

func TestAppendTextValid(t *testing.T) {
	var testdata []ValidExample
	err := yaml.Unmarshal(validEncodingYML, &testdata)
	require.NoError(t, err)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			tid := typeid.MustParse(td.Tid)

			// Test AppendText with nil slice (equivalent to MarshalText)
			result, err := tid.AppendText(nil)
			assert.NoError(t, err)
			assert.Equal(t, td.Tid, string(result))

			// Test AppendText with existing data
			prefix := []byte("prefix:")
			result, err = tid.AppendText(prefix)
			assert.NoError(t, err)
			assert.Equal(t, "prefix:"+td.Tid, string(result))

			// Verify that MarshalText and AppendText(nil) are semantically identical
			marshaled, err := tid.MarshalText()
			assert.NoError(t, err)
			appended, err := tid.AppendText(nil)
			assert.NoError(t, err)
			assert.Equal(t, marshaled, appended, "MarshalText() should be identical to AppendText(nil)")

			// Test that the original slice is not modified
			original := []byte("original")
			originalCopy := make([]byte, len(original))
			copy(originalCopy, original)
			result, err = tid.AppendText(original)
			assert.NoError(t, err)
			assert.Equal(t, originalCopy, original, "original slice should not be modified")
			assert.Equal(t, "original"+td.Tid, string(result))
		})
	}
}

func TestJSONInvalid(t *testing.T) {
	var testdata []InvalidExample
	err := yaml.Unmarshal(invalidEncodingYML, &testdata)
	require.NoError(t, err)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			// Test UnmarshalText with invalid TypeID strings
			var decoded typeid.TypeID
			invalidJSON := `"` + td.Tid + `"`
			err := json.Unmarshal([]byte(invalidJSON), &decoded)
			assert.Error(t, err, "JSON unmarshal should fail for invalid typeid: %s", td.Tid)
		})
	}
}

// Mock structs for testing JSON encoding behavior
type MockWithoutOmitZero struct {
	ID typeid.TypeID `json:"id"`
}

type MockWithOmitZero struct {
	ID typeid.TypeID `json:"id,omitzero"`
}

func TestJSONOmitZero(t *testing.T) {
	testCases := []struct {
		name            string
		typeID          typeid.TypeID
		expectedWithout string // expected JSON without omitzero
		expectedWith    string // expected JSON with omitzero
		description     string
	}{
		{
			name:            "zero TypeID",
			typeID:          typeid.TypeID{},
			expectedWithout: `{"id":"00000000000000000000000000"}`,
			expectedWith:    `{}`,
			description:     "empty TypeID struct should omit with omitzero tag",
		},
		{
			name:            "constructed zero ID",
			typeID:          typeid.MustParse("00000000000000000000000000"),
			expectedWithout: `{"id":"00000000000000000000000000"}`,
			expectedWith:    `{}`,
			description:     "constructed zero ID should omit with omitzero tag",
		},
		{
			name:            "prefixed zero ID",
			typeID:          typeid.MustParse("user_00000000000000000000000000"),
			expectedWithout: `{"id":"user_00000000000000000000000000"}`,
			expectedWith:    `{"id":"user_00000000000000000000000000"}`,
			description:     "prefixed zero ID should not omit because IsZero() returns false",
		},
		{
			name:            "non-zero ID",
			typeID:          typeid.MustParse("prefix_01h455vb4pex5vsknk084sn02q"),
			expectedWithout: `{"id":"prefix_01h455vb4pex5vsknk084sn02q"}`,
			expectedWith:    `{"id":"prefix_01h455vb4pex5vsknk084sn02q"}`,
			description:     "non-zero ID should always be included",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test without omitzero
			t.Run("without omitzero", func(t *testing.T) {
				mock := MockWithoutOmitZero{ID: tc.typeID}
				encoded, err := json.Marshal(mock)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedWithout, string(encoded), tc.description)

				// Test roundtrip
				var decoded MockWithoutOmitZero
				err = json.Unmarshal(encoded, &decoded)
				assert.NoError(t, err)
				assert.Equal(t, mock, decoded)
			})

			// Test with omitzero
			t.Run("with omitzero", func(t *testing.T) {
				mock := MockWithOmitZero{ID: tc.typeID}
				encoded, err := json.Marshal(mock)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedWith, string(encoded), tc.description)

				// Test roundtrip - for omitted fields, unmarshal from empty object
				var decoded MockWithOmitZero
				if tc.expectedWith == `{}` {
					err = json.Unmarshal([]byte(`{}`), &decoded)
				} else {
					err = json.Unmarshal(encoded, &decoded)
				}
				assert.NoError(t, err)
				assert.Equal(t, mock, decoded)
			})
		})
	}
}
