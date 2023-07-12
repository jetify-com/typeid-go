package typeid_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.jetpack.io/typeid"
)

func TestJSON(t *testing.T) {
	str := "prefix_00041061050r3gg28a1c60t3gf"
	tid := typeid.Must(typeid.FromString(str))

	encoded, err := json.Marshal(tid)
	assert.NoError(t, err)
	assert.Equal(t, `"`+str+`"`, string(encoded))

	var decoded typeid.TypeID
	err = json.Unmarshal(encoded, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, tid, decoded)
	assert.Equal(t, str, decoded.String())
}
