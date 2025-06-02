package base32

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	encoder := base32.NewEncoding("0123456789abcdefghjkmnpqrstvwxyz")

	for i := 0; i < 1000; i++ {
		// Generate 16 random bytes
		data := make([]byte, 16)
		_, err := rand.Read(data)
		assert.NoError(t, err)

		// Encode them using our library, and encode them using go's standard library:
		actual := EncodeToString([16]byte(data))

		// The standard base32 library decodes in groups of 5 bytes, otherwise it needs
		// to pad, by default it pads at the end of the byte array, but to match our
		// encoding we need to pad in the front.
		// Pad manually, and then remove the extra 000000 from the resulting string.
		padded := append([]byte{0x00, 0x00, 0x00, 0x00}, data...)
		expected := encoder.EncodeToString(padded)[6:]

		// They should be equal
		assert.Equal(t, expected, actual)

		// Decoding again should yield the original result:
		decoded, err := DecodeString(actual)
		assert.NoError(t, err)
		assert.Equal(t, data, decoded[:])
	}
}

// TestEncodeToBuffer tests the zero-allocation encoding to a pre-allocated buffer
func TestEncode(t *testing.T) {
	encoder := base32.NewEncoding("0123456789abcdefghjkmnpqrstvwxyz")

	for i := 0; i < 1000; i++ {
		// Generate 16 random bytes
		data := make([]byte, 16)
		_, err := rand.Read(data)
		assert.NoError(t, err)

		// Test with exact size buffer
		dst := make([]byte, 26)
		Encode(dst, [16]byte(data))

		// Compare with stdlib
		padded := append([]byte{0x00, 0x00, 0x00, 0x00}, data...)
		expected := encoder.EncodeToString(padded)[6:]
		assert.Equal(t, expected, string(dst))

		// Test with larger buffer
		largeDst := make([]byte, 50)
		Encode(largeDst, [16]byte(data))
		assert.Equal(t, expected, string(largeDst[:26]))

		// Test Decode method
		decodeDst := make([]byte, 16)
		n, err := Decode(decodeDst, dst)
		assert.NoError(t, err)
		assert.Equal(t, 16, n)
		assert.Equal(t, data, decodeDst)
	}
}

// TestAppendEncode tests the append-style encoding
func TestAppendEncode(t *testing.T) {
	encoder := base32.NewEncoding("0123456789abcdefghjkmnpqrstvwxyz")

	// Test appending to empty slice
	data1 := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	result := AppendEncode(nil, data1)
	assert.Len(t, result, 26)

	// Compare with stdlib
	padded := append([]byte{0x00, 0x00, 0x00, 0x00}, data1[:]...)
	expected := encoder.EncodeToString(padded)[6:]
	assert.Equal(t, expected, string(result))

	// Test appending to existing slice
	data2 := [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	result = AppendEncode(result, data2)
	assert.Len(t, result, 52)

	// Verify both encodings are correct
	assert.Equal(t, expected, string(result[:26]))

	padded2 := append([]byte{0x00, 0x00, 0x00, 0x00}, data2[:]...)
	expected2 := encoder.EncodeToString(padded2)[6:]
	assert.Equal(t, expected2, string(result[26:]))
}

// TestAppendDecode tests the append-style decoding
func TestAppendDecode(t *testing.T) {
	data1 := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	data2 := [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	encoded1 := EncodeToString(data1)
	encoded2 := EncodeToString(data2)

	// Test appending to empty slice
	result, err := AppendDecode(nil, []byte(encoded1))
	assert.NoError(t, err)
	assert.Equal(t, data1[:], result)

	// Test appending to existing slice
	result, err = AppendDecode(result, []byte(encoded2))
	assert.NoError(t, err)
	assert.Len(t, result, 32)
	assert.Equal(t, data1[:], result[:16])
	assert.Equal(t, data2[:], result[16:])
}

// TestEncodePanic tests that Encode panics with insufficient buffer
func TestEncodePanic(t *testing.T) {
	data := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	smallBuffer := make([]byte, 25) // One byte too small

	assert.Panics(t, func() {
		Encode(smallBuffer, data)
	})
}

// TestCorruptInputError tests that CorruptInputError reports specific positions
func TestCorruptInputError(t *testing.T) {
	// Test with invalid character at position 5
	invalidInput := "01234!6789abcdefghjkmnpqrs"
	_, err := DecodeString(invalidInput)
	assert.Error(t, err)

	var corruptErr CorruptInputError
	assert.ErrorAs(t, err, &corruptErr)
	assert.Equal(t, CorruptInputError(5), corruptErr)
	assert.Contains(t, err.Error(), "illegal base32 data at offset 5")

	// Test with invalid character at position 0
	invalidInput0 := "!123456789abcdefghjkmnpqrs"
	_, err = DecodeString(invalidInput0)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &corruptErr)
	assert.Equal(t, CorruptInputError(0), corruptErr)
	assert.Contains(t, err.Error(), "illegal base32 data at offset 0")

	// Test with wrong length
	shortInput := "01234"
	_, err = DecodeString(shortInput)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &corruptErr)
	assert.Equal(t, CorruptInputError(26), corruptErr) // Length error reports position 26
	assert.Contains(t, err.Error(), "illegal base32 data at offset 26")
}

// TestKnownVectors tests encoding/decoding with known test vectors from the spec
func TestKnownVectors(t *testing.T) {
	testCases := []struct {
		name     string
		input    [16]byte
		expected string
	}{
		{
			name:     "nil",
			input:    [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected: "00000000000000000000000000",
		},
		{
			name:     "one",
			input:    [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			expected: "00000000000000000000000001",
		},
		{
			name:     "ten",
			input:    [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a},
			expected: "0000000000000000000000000a",
		},
		{
			name:     "sixteen",
			input:    [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10},
			expected: "0000000000000000000000000g",
		},
		{
			name:     "thirty-two",
			input:    [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20},
			expected: "00000000000000000000000010",
		},
		{
			name:     "max-valid",
			input:    [16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			expected: "7zzzzzzzzzzzzzzzzzzzzzzzzz",
		},
		{
			name:     "valid-alphabet",
			input:    [16]byte{0x01, 0x10, 0xc8, 0x53, 0x1d, 0x09, 0x52, 0xd8, 0xd7, 0x3e, 0x11, 0x94, 0xe9, 0x5b, 0x5f, 0x19},
			expected: "0123456789abcdefghjkmnpqrs",
		},
		{
			name:     "valid-uuidv7",
			input:    [16]byte{0x01, 0x89, 0x0a, 0x5d, 0xac, 0x96, 0x77, 0x4b, 0xbc, 0xce, 0xb3, 0x02, 0x09, 0x9a, 0x80, 0x57},
			expected: "01h455vb4pex5vsknk084sn02q",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test encoding
			encoded := EncodeToString(tc.input)
			assert.Equal(t, tc.expected, encoded, "encoding should match expected value")

			// Test decoding
			decoded, err := DecodeString(tc.expected)
			assert.NoError(t, err, "decoding should not return an error")
			assert.Equal(t, tc.input[:], decoded[:], "decoded value should match original input")

			// Test round trip
			var decodedArray [16]byte
			copy(decodedArray[:], decoded)
			roundTrip := EncodeToString(decodedArray)
			assert.Equal(t, tc.expected, roundTrip, "round trip should preserve the encoding")
		})
	}
}

// TestBufferSizes tests various buffer size scenarios
func TestBufferSizes(t *testing.T) {
	data := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	expected := EncodeToString(data)

	t.Run("exact size buffer", func(t *testing.T) {
		dst := make([]byte, 26)
		Encode(dst, data)
		assert.Equal(t, expected, string(dst))
	})

	t.Run("oversized buffer", func(t *testing.T) {
		dst := make([]byte, 100)
		Encode(dst, data)
		assert.Equal(t, expected, string(dst[:26]))
		// Verify the rest of the buffer is unchanged (should be zeros)
		assert.Equal(t, make([]byte, 74), dst[26:])
	})

	t.Run("undersized buffer panics", func(t *testing.T) {
		dst := make([]byte, 25)
		assert.Panics(t, func() {
			Encode(dst, data)
		}, "encoding to undersized buffer should panic")
	})
}

// TestInvalidInputs tests various invalid input scenarios
func TestInvalidInputs(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError bool
		errorContains string
	}{
		{
			name:          "invalid character",
			input:         "01234!6789abcdefghjkmnpqrs", // contains invalid char '!'
			expectedError: true,
			errorContains: "illegal base32 data",
		},
		{
			name:          "wrong length - too short",
			input:         "0123456789",
			expectedError: true,
			errorContains: "illegal base32 data",
		},
		{
			name:          "wrong length - too long",
			input:         "012345678901234567890123456789",
			expectedError: true,
			errorContains: "illegal base32 data",
		},
		{
			name:          "empty string",
			input:         "",
			expectedError: true,
			errorContains: "illegal base32 data",
		},
		{
			name:          "valid input",
			input:         "00041061050r3gg28a1c60t3gf",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeString(tt.input)

			if tt.expectedError {
				assert.Error(t, err, "should return an error for invalid input")
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, "error message should contain expected text")
				}
			} else {
				assert.NoError(t, err, "should not return an error for valid input")
			}
		})
	}
}

// TestInvalidCharacterAtEachPosition tests that CorruptInputError reports correct positions
// for invalid characters at each of the 26 positions in the input string
func TestInvalidCharacterAtEachPosition(t *testing.T) {
	// Start with a valid base32 string
	validInput := "01h455vb4pex5vsknk084sn02q"

	for pos := 0; pos < 26; pos++ {
		t.Run(fmt.Sprintf("position_%d", pos), func(t *testing.T) {
			// Create invalid input by replacing character at position with '!'
			invalidInput := []byte(validInput)
			invalidInput[pos] = '!'

			_, err := DecodeString(string(invalidInput))
			assert.Error(t, err, "should return error for invalid character")

			var corruptErr CorruptInputError
			assert.ErrorAs(t, err, &corruptErr, "should be CorruptInputError")
			assert.Equal(t, CorruptInputError(pos), corruptErr, "should report correct position")
			assert.Contains(t, err.Error(), fmt.Sprintf("illegal base32 data at offset %d", pos))
		})
	}
}

// TestAppendDecodeErrorHandling tests error handling in AppendDecode
func TestAppendDecodeErrorHandling(t *testing.T) {
	// Test with invalid input that will cause decode to fail
	invalidInput := "01h455vb4pex5vsknk084sn02!" // Invalid character at end

	// Add some initial data to the destination
	initialData := []byte("initial")
	result, err := AppendDecode(initialData, []byte(invalidInput))

	assert.Error(t, err, "should return error for invalid input")
	assert.Equal(t, initialData, result, "should return original data when decode fails")

	var corruptErr CorruptInputError
	assert.ErrorAs(t, err, &corruptErr, "should be CorruptInputError")
	assert.Equal(t, CorruptInputError(25), corruptErr, "should report correct position")
}
