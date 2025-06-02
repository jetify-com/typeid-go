package base32

// Encoding and Decoding code based on the go implementation of ulid
// found at: https://github.com/oklog/ulid
// (Copyright 2016 The Oklog Authors)
// Modifications made available under the same license as the original

import (
	"slices"
	"strconv"
)

const alphabet = "0123456789abcdefghjkmnpqrstvwxyz"

// CorruptInputError is returned when invalid base32 data is encountered.
type CorruptInputError int

func (e CorruptInputError) Error() string {
	return "illegal base32 data at offset " + strconv.FormatInt(int64(e), 10)
}

// Encode encodes src using the base32 alphabet,
// writing exactly 26 bytes to dst.
// The caller must ensure that dst is large enough to hold all the encoded data.
func Encode(dst []byte, src [16]byte) {
	// Optimized unrolled loop ahead.

	// 10 byte timestamp
	dst[0] = alphabet[(src[0]&224)>>5]
	dst[1] = alphabet[src[0]&31]
	dst[2] = alphabet[(src[1]&248)>>3]
	dst[3] = alphabet[((src[1]&7)<<2)|((src[2]&192)>>6)]
	dst[4] = alphabet[(src[2]&62)>>1]
	dst[5] = alphabet[((src[2]&1)<<4)|((src[3]&240)>>4)]
	dst[6] = alphabet[((src[3]&15)<<1)|((src[4]&128)>>7)]
	dst[7] = alphabet[(src[4]&124)>>2]
	dst[8] = alphabet[((src[4]&3)<<3)|((src[5]&224)>>5)]
	dst[9] = alphabet[src[5]&31]

	// 16 bytes of entropy
	dst[10] = alphabet[(src[6]&248)>>3]
	dst[11] = alphabet[((src[6]&7)<<2)|((src[7]&192)>>6)]
	dst[12] = alphabet[(src[7]&62)>>1]
	dst[13] = alphabet[((src[7]&1)<<4)|((src[8]&240)>>4)]
	dst[14] = alphabet[((src[8]&15)<<1)|((src[9]&128)>>7)]
	dst[15] = alphabet[(src[9]&124)>>2]
	dst[16] = alphabet[((src[9]&3)<<3)|((src[10]&224)>>5)]
	dst[17] = alphabet[src[10]&31]
	dst[18] = alphabet[(src[11]&248)>>3]
	dst[19] = alphabet[((src[11]&7)<<2)|((src[12]&192)>>6)]
	dst[20] = alphabet[(src[12]&62)>>1]
	dst[21] = alphabet[((src[12]&1)<<4)|((src[13]&240)>>4)]
	dst[22] = alphabet[((src[13]&15)<<1)|((src[14]&128)>>7)]
	dst[23] = alphabet[(src[14]&124)>>2]
	dst[24] = alphabet[((src[14]&3)<<3)|((src[15]&224)>>5)]
	dst[25] = alphabet[src[15]&31]
}

// AppendEncode appends the base32 encoded src to dst and returns the extended buffer.
// This is efficient for building up encoded data without repeated allocations.
func AppendEncode(dst []byte, src [16]byte) []byte {
	dst = slices.Grow(dst, 26)
	start := len(dst)
	dst = dst[:start+26]
	Encode(dst[start:], src)
	return dst
}

// EncodeToString returns the base32 encoding of src.
func EncodeToString(src [16]byte) string {
	dst := make([]byte, 26)
	Encode(dst, src)
	return string(dst)
}

// Byte to index table for O(1) lookups when unmarshaling.
// We use 0xFF as sentinel value for invalid indexes.
var dec = [...]byte{
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x01,
	0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x0A, 0x0B, 0x0C,
	0x0D, 0x0E, 0x0F, 0x10, 0x11, 0xFF, 0x12, 0x13, 0xFF, 0x14,
	0x15, 0xFF, 0x16, 0x17, 0x18, 0x19, 0x1A, 0xFF, 0x1B, 0x1C,
	0x1D, 0x1E, 0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
}

// ValidateBytes validates that src is exactly 26 bytes and contains
// only valid base32 characters. Returns CorruptInputError with position
// if invalid.
func ValidateBytes(src []byte) error {
	if len(src) != 26 {
		return CorruptInputError(26)
	}

	// Check if all the characters are part of the expected base32 character set.
	// Optimized unrolled validation with specific error positions
	if dec[src[0]] == 0xFF {
		return CorruptInputError(0)
	}
	if dec[src[1]] == 0xFF {
		return CorruptInputError(1)
	}
	if dec[src[2]] == 0xFF {
		return CorruptInputError(2)
	}
	if dec[src[3]] == 0xFF {
		return CorruptInputError(3)
	}
	if dec[src[4]] == 0xFF {
		return CorruptInputError(4)
	}
	if dec[src[5]] == 0xFF {
		return CorruptInputError(5)
	}
	if dec[src[6]] == 0xFF {
		return CorruptInputError(6)
	}
	if dec[src[7]] == 0xFF {
		return CorruptInputError(7)
	}
	if dec[src[8]] == 0xFF {
		return CorruptInputError(8)
	}
	if dec[src[9]] == 0xFF {
		return CorruptInputError(9)
	}
	if dec[src[10]] == 0xFF {
		return CorruptInputError(10)
	}
	if dec[src[11]] == 0xFF {
		return CorruptInputError(11)
	}
	if dec[src[12]] == 0xFF {
		return CorruptInputError(12)
	}
	if dec[src[13]] == 0xFF {
		return CorruptInputError(13)
	}
	if dec[src[14]] == 0xFF {
		return CorruptInputError(14)
	}
	if dec[src[15]] == 0xFF {
		return CorruptInputError(15)
	}
	if dec[src[16]] == 0xFF {
		return CorruptInputError(16)
	}
	if dec[src[17]] == 0xFF {
		return CorruptInputError(17)
	}
	if dec[src[18]] == 0xFF {
		return CorruptInputError(18)
	}
	if dec[src[19]] == 0xFF {
		return CorruptInputError(19)
	}
	if dec[src[20]] == 0xFF {
		return CorruptInputError(20)
	}
	if dec[src[21]] == 0xFF {
		return CorruptInputError(21)
	}
	if dec[src[22]] == 0xFF {
		return CorruptInputError(22)
	}
	if dec[src[23]] == 0xFF {
		return CorruptInputError(23)
	}
	if dec[src[24]] == 0xFF {
		return CorruptInputError(24)
	}
	if dec[src[25]] == 0xFF {
		return CorruptInputError(25)
	}
	return nil
}

// ValidateString validates that s is exactly 26 characters and contains
// only valid base32 characters. Returns CorruptInputError with position
// if invalid.
func ValidateString(s string) error {
	return ValidateBytes([]byte(s))
}

// decode is the core decoding logic, used by other decode methods.
// It validates the input and decodes it into dst.
func decode(dst, src []byte) (n int, err error) {
	// Validate the input using the shared validation logic
	if err := ValidateBytes(src); err != nil {
		return 0, err
	}

	// 6 bytes timestamp (48 bits)
	dst[0] = (dec[src[0]] << 5) | dec[src[1]]
	dst[1] = (dec[src[2]] << 3) | (dec[src[3]] >> 2)
	dst[2] = (dec[src[3]] << 6) | (dec[src[4]] << 1) | (dec[src[5]] >> 4)
	dst[3] = (dec[src[5]] << 4) | (dec[src[6]] >> 1)
	dst[4] = (dec[src[6]] << 7) | (dec[src[7]] << 2) | (dec[src[8]] >> 3)
	dst[5] = (dec[src[8]] << 5) | dec[src[9]]

	// 10 bytes of entropy (80 bits)
	dst[6] = (dec[src[10]] << 3) | (dec[src[11]] >> 2) // First 4 bits are the version
	dst[7] = (dec[src[11]] << 6) | (dec[src[12]] << 1) | (dec[src[13]] >> 4)
	dst[8] = (dec[src[13]] << 4) | (dec[src[14]] >> 1) // First 2 bits are the variant
	dst[9] = (dec[src[14]] << 7) | (dec[src[15]] << 2) | (dec[src[16]] >> 3)
	dst[10] = (dec[src[16]] << 5) | dec[src[17]]
	dst[11] = (dec[src[18]] << 3) | dec[src[19]]>>2
	dst[12] = (dec[src[19]] << 6) | (dec[src[20]] << 1) | (dec[src[21]] >> 4)
	dst[13] = (dec[src[21]] << 4) | (dec[src[22]] >> 1)
	dst[14] = (dec[src[22]] << 7) | (dec[src[23]] << 2) | (dec[src[24]] >> 3)
	dst[15] = (dec[src[24]] << 5) | dec[src[25]]

	return 16, nil
}

// Decode decodes src using the base32 alphabet into dst.
// It writes exactly 16 bytes to dst and returns the number of bytes written.
// The caller must ensure that dst is large enough to hold all the decoded data.
// If src contains invalid base32 data, it will return CorruptInputError.
func Decode(dst, src []byte) (n int, err error) {
	return decode(dst, src)
}

// AppendDecode appends the base32 decoded src to dst and returns the extended buffer.
// If the input is malformed, it returns the original dst and an error.
func AppendDecode(dst, src []byte) ([]byte, error) {
	dst = slices.Grow(dst, 16)
	start := len(dst)
	dst = dst[:start+16]

	n, err := decode(dst[start:], src)
	if err != nil {
		return dst[:start], err
	}

	return dst[:start+n], nil
}

// DecodeString returns the bytes represented by the base32 string s.
// If the input is malformed, it returns nil and CorruptInputError.
func DecodeString(s string) ([]byte, error) {
	dst := make([]byte, 16)
	n, err := decode(dst, []byte(s))
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}
