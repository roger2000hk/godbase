package godbase

import (
	"encoding/binary"
	"github.com/wayn3h0/go-uuid"
)

type UId uuid.UUID

var ByteOrder = binary.BigEndian

func NewUId() UId {
	res, err := uuid.NewRandom()

	if err != nil {
		panic(err)
	}

	return UId(res)
}
