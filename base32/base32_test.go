package base32

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInverse(t *testing.T) {
	data := [16]byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	}
	encoded := Encode(data)
	assert.Equal(t, "00041061050r3gg28a1c60t3gf", encoded)
	decoded, err := Decode(encoded)
	assert.NoError(t, err)
	for i := 0; i < 16; i++ {
		assert.Equal(t, data[i], decoded[i])
	}
}
