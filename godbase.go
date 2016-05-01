package godbase

import (
	"encoding/binary"
	"github.com/wayn3h0/go-uuid"
	"hash"
	"hash/fnv"
	"io"
	"time"
)

type Col interface {
	Def
	AsKey(Rec, interface{}) Key
	CloneVal(interface{}) interface{}
	Decode(interface{}) interface{}
	Encode(interface{}) interface{}
	Eq(interface{}, interface{}) bool
	Hash(Rec, interface{}, hash.Hash64)
	Read(Rec, ValSize, io.Reader) (interface{}, error)
	Type() ColType
	Write(Rec, interface{}, io.Writer) error
}

type ColType interface {
	AsKey(Rec, interface{}) Key
	CloneVal(interface{}) interface{}
	Decode(interface{}) interface{}
	Encode(interface{}) interface{}
	Eq(interface{}, interface{}) bool
	Hash(Rec, interface{}, hash.Hash64)
	Name() string
	Read(Rec, ValSize, io.Reader) (interface{}, error)
	Write(Rec, interface{}, io.Writer) error
}

type Cx interface {
	InitRec(Rec) Rec
	InitRecId(Rec, UId) Rec
}

type Def interface {
	Key
	Name() string
}

type Idx interface {
	Def
	Delete(Iter, Rec) error
	Find(start Iter, key Key, val interface{}) (Iter, bool)
	Insert(Iter, Rec) (Iter, error)
	Load(Rec) (Rec, error)
	Key(...interface{}) Key
	RecKey(r Rec) Key
}

type Map interface {
	// Clears all elems from map. Deallocates nodes for maps that use allocators.
	Clear()

	// Cuts elems from start to end for which fn returns non nil key into new set;
	// start, end & fn are all optional. When fn is specified, the returned key/val replaces
	// the original; except for maps with embedded nodes, where the returned val replaces the
	// entire node. No safety checks are provided; if you mess up the ordering, you're on your
	// own. Circular cuts, with start/end on opposite sides of root; are supported. 
	// Returns a cut from the start slot for hash maps.

	Cut(start, end Iter, fn KVMapFn) Map

	// Deletes elems from start to end, matching key/val;
	// start, end, key & val are all optional, nil means all elems. Specifying 
	// iters for hash maps only works within the same slot. Circular deletes,
	// with start/end on opposite sides of root; are supported. Deallocates nodes
	// for maps that use allocators. Returns an iter to next 
	// elem and number of deleted elems.

	Delete(start, end Iter, key Key, val interface{}) (Iter, int)

	// Returns iter for first elem after start matching key and ok;
	// start & val are optional, specifying a start iter for hash maps only works within the 
	// same slot.

	Find(start Iter, key Key, val interface{}) (Iter, bool)
	
	// Returns iter to first elem; not supported by hash maps
	First() Iter

	// Returns val for key and ok
	Get(key Key) (interface{}, bool)

	// Inserts key/val into map after start;
	// start & val are both optional, dup checks can be disabled by setting allowMulti to false. 
	// Returns iter to inserted val & true on success, or iter to existing val & false on dup. 
	// Specifying a start iter for hash maps only works within the same slot.

	Insert(start Iter, key Key, val interface{}, allowMulti bool) (Iter, bool)

	// Returns a new, empty map of the same type as the receiver
	New() Map

	// Returns the number of elems in map
	Len() int64

	// Inserts/updates key to val and returns true on insert
	Set(key Key, val interface{}) bool

	// Returns string repr for printing
	String() string
	
	// Calls fn with successive elems until it returns false; returns false on early exit
	While(KVTestFn) bool
}

type Rec interface {
	Clear()
	Clone() Rec
	Delete(Col) bool
	Eq(Rec) bool
	Find(Col) (interface{}, bool)
	Get(Col) interface{}
	Id() UId
	Iter() Iter
	Len() int
	New() Rec
	Set(Col, interface{}) interface{}
}

type Tbl interface {
	Def
	AddCol(Col) Col
	Clear()
	Col(string) Col
	Cols() Iter
	Dump(io.Writer) error
	Len() int64
	Load(Cx, Rec) (Rec, error)
	Reset(Rec) (Rec, error)
	Read(Rec, io.Reader) (Rec, error)
	Revision(Rec) int64
	OnDelete() *Evt
	OnLoad() *Evt
	OnUpsert() *Evt
	Slurp(Cx, io.Reader) error
	Upsert(Cx, Rec) (Rec, error)
	UpsertedAt(Rec) time.Time
	Write(Rec, io.Writer) error
}

type KVMapFn func (Key, interface{}) (Key, interface{})
type KVTestFn func (Key, interface{}) bool
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
