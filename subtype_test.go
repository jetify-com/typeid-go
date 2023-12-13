package typeid_test

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"go.jetpack.io/typeid"
)

func ExampleNew() {
	tid := typeid.Must(typeid.New[AccountID]())
	fmt.Println("Prefix:", tid.Prefix())
	// Output:
	// Prefix: account
}

func ExampleFromSuffix() {
	tid := typeid.Must(typeid.FromSuffix[UserID]("00041061050r3gg28a1c60t3gf"))
	fmt.Printf("Prefix: %s\nSuffix: %s\n", tid.Prefix(), tid.Suffix())
	// Output:
	// Prefix: user
	// Suffix: 00041061050r3gg28a1c60t3gf
}

func TestSubtypeConstructors(t *testing.T) {
	// These constructors should work for a subtype:
	_, err := typeid.New[AccountID]()
	assert.NoError(t, err)
	_, err = typeid.FromSuffix[AccountID]("00041061050r3gg28a1c60t3gf")
	assert.NoError(t, err)

	// But error on TypeID[typeid.Any]:
	_, err = typeid.New[typeid.AnyID]()
	assert.Error(t, err)
	_, err = typeid.FromSuffix[typeid.AnyID]("00041061050r3gg28a1c60t3gf")
	assert.Error(t, err)
}

func TestSubtypeNil(t *testing.T) {
	var emptyUser UserID
	nilUser, err := typeid.Parse[UserID]("user_00000000000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, nilUser, emptyUser)
	assert.Equal(t, nilUser.String(), emptyUser.String())
	assert.Equal(t, nilUser.Prefix(), emptyUser.Prefix())
	assert.Equal(t, nilUser.UUID(), emptyUser.UUID())
	assert.Equal(t, nilUser.UUIDBytes(), emptyUser.UUIDBytes())
	assert.Equal(t, "user_00000000000000000000000000", nilUser.String())
	assert.Equal(t, "user", nilUser.Prefix())

	parsed, err := typeid.FromString("user_00000000000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, "user_00000000000000000000000000", parsed.String())
	assert.Equal(t, "user", parsed.Prefix())
	assert.Equal(t, "00000000000000000000000000", parsed.Suffix())
}

func TestParse(t *testing.T) {
	// Generate a bunch of random UserIDs. We should be able to parse them
	// using the correct type, but not an incorrect one.
	for i := 0; i < 1000; i++ {
		tid := typeid.Must(typeid.New[UserID]())
		// They parse as UserID
		parsed, err := typeid.Parse[UserID](tid.String())
		if err != nil {
			t.Error(err)
		}
		if tid != parsed {
			t.Errorf("Expected %s, got %s", tid, parsed)
		}

		// They also parse as a generic TypeID
		_, err = typeid.FromString(tid.String())
		if err != nil {
			t.Error(err)
		}

		// But not as an AccountID
		_, err = typeid.Parse[AccountID](tid.String())
		assert.Error(t, err)
	}
}

func TestFromUUID(t *testing.T) {
	uid, err := uuid.NewV7()
	assert.NoError(t, err)
	id, err := typeid.FromUUID[UserID](uid.String())
	assert.NoError(t, err)
	assert.Equal(t, uid.String(), id.UUID())
}

func TestFromUUIDBytes(t *testing.T) {
	uid, err := uuid.NewV7()
	assert.NoError(t, err)
	id, err := typeid.FromUUIDBytes[UserID](uid.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, uid.Bytes(), id.UUIDBytes())
}
