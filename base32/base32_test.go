package base32

import (
	"crypto/rand"
	"encoding/base32"
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
		actual := Encode([16]byte(data))

		// The standard base32 library decodes in groups of 5 bytes, otherwise it needs
		// to pad, by default it pads at the end of the byte array, but to match our
		// encoding we need to pad in the front.
		// Pad manually, and then remove the extra 000000 from the resulting string.
		padded := append([]byte{0x00, 0x00, 0x00, 0x00}, data...)
		expected := encoder.EncodeToString(padded)[6:]

		// They should be equal
		assert.Equal(t, expected, actual)

		// Decoding again should yield the original result:
		decoded, err := Decode(actual)
		assert.NoError(t, err)
		for i := 0; i < 16; i++ {
			assert.Equal(t, data[i], decoded[i])
		}
	}
}
