package godbase

import (
	"encoding/binary"
	"github.com/wayn3h0/go-uuid"
	"hash"
	"hash/fnv"
	"io"
)

type Cx interface {
	InitRec(Rec) Rec
	InitRecId(Rec, UId) Rec
}

type Def interface {
	Key
	Name() string
}

type HashFn func (Key) uint64
type NameSize uint8
type UId uuid.UUID
type ValSize uint32

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

func ReadUId(r io.Reader) (interface{}, error) {
	var v UId

	if _, err := io.ReadFull(r, v[:]); err != nil {
		return nil, err
	}

	return v, nil
}

func Write(ptr interface{}, w io.Writer) error {
	return binary.Write(w, ByteOrder, ptr)
}

type UIdHash struct {
	imp hash.Hash64
}

func NewUIdHash() *UIdHash {
	return new(UIdHash).Init()
}

func (h *UIdHash) Hash(id UId) uint64 {
	h.imp.Reset()
	h.imp.Write(id[:])
	return h.imp.Sum64()
}

func (h *UIdHash) Init() *UIdHash {
	h.imp = fnv.New64()
	return h
}
