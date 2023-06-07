package typeid

import (
	"github.com/gofrs/uuid"
	"go.jetpack.io/typeid/base32"
)

type TypeID struct {
	prefix string
	suffix string
}

var Nil = TypeID{
	prefix: "",
	suffix: "00000000000000000000000000",
}

func New(prefix string) (TypeID, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return Nil, err
	}
	tid := TypeID{
		prefix: prefix,
		suffix: base32.Encode(uid),
	}
	return tid, nil
}

func (tid TypeID) Type() string {
	return tid.prefix
}

func (tid TypeID) Suffix() string {
	return tid.suffix
}

func (tid TypeID) String() string {
	if tid.prefix == "" {
		return tid.suffix
	}
	return tid.prefix + "_" + tid.Suffix()
}

func (tid TypeID) UUIDBytes() []byte {
	b, err := base32.Decode(tid.suffix)
	if err != nil {
		panic(err)
	}
	return b
}

func (tid TypeID) UUIDString() string {
	return uuid.FromBytesOrNil(tid.UUIDBytes()).String()
}

func Must(tid TypeID, err error) TypeID {
	if err != nil {
		panic(err)
	}
	return tid
}
