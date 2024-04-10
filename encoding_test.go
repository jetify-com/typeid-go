package typeid_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.jetify.com/typeid"
)

func TestJSON(t *testing.T) {
	str := "prefix_00041061050r3gg28a1c60t3gf"
	tid := typeid.Must(typeid.FromString(str))

	encoded, err := json.Marshal(tid)
	assert.NoError(t, err)
	assert.Equal(t, `"`+str+`"`, string(encoded))

	var decoded typeid.AnyID
	err = json.Unmarshal(encoded, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, tid, decoded)
	assert.Equal(t, str, decoded.String())
}

func TestJSON_Subtype(t *testing.T) {
	str := "user_00041061050r3gg28a1c60t3gf"
	tid := typeid.Must(typeid.Parse[UserID](str))

	encoded, err := json.Marshal(tid)
	assert.NoError(t, err)
	assert.Equal(t, `"`+str+`"`, string(encoded))

	var decoded UserID
	err = json.Unmarshal(encoded, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, tid, decoded)
	assert.Equal(t, str, decoded.String())

	var wrongType AccountID
	err = json.Unmarshal(encoded, &wrongType)
	assert.Error(t, err)
}
