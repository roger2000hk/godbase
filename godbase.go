package godbase

import (
	"encoding/binary"
	"github.com/wayn3h0/go-uuid"
	"io"
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


func (id UId) String() string {
	return uuid.UUID(id).String()
}

func Read(ptr interface{}, r io.Reader) error {
	if err := binary.Read(r, ByteOrder, ptr); err != nil {
		return err
	}

	return nil
}

func Write(ptr interface{}, w io.Writer) error {
	return binary.Write(w, ByteOrder, ptr)
}
