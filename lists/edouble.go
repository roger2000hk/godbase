package lists

import (
	//"fmt"
	"unsafe"
)

type EDouble struct {
	next, prev *EDouble
}

func EmptyEDouble() *EDouble {
	return new(EDouble).Init()
}

func (l *EDouble) Del() {
	l.prev.next = l.next
	l.next.prev = l.prev
	l.next, l.prev = l, l
}

func (l *EDouble) Init() *EDouble {
	l.next, l.prev = l, l
	return l
}

func (l *EDouble) InsAfter(e *EDouble) *EDouble {
	e.next, e.prev = l.next, l
	l.next.prev, l.next = e, e
	return e
}

func (l *EDouble) InsBefore(e *EDouble) *EDouble {
	e.next, e.prev = l, l.prev
	l.prev.next, l.prev = e, e
	return e
}

func (l *EDouble) Next() *EDouble {
	return l.next
}

func (l *EDouble) Prev() *EDouble {
	return l.prev
}

func (l *EDouble) Ptr(fld uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(l)) - fld)
}

func (l *EDouble) Val() interface{} {
	return l
}
