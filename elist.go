package gofbls

import (
	//"fmt"
	"unsafe"
)

type EList struct {
	next, prev *EList
}

func EmptyEList() *EList {
	return new(EList).Init()
}

func (l *EList) Del() {
	l.prev.next = l.next
	l.next.prev = l.prev
	l.next, l.prev = l, l
}

func (l *EList) Init() *EList {
	l.next, l.prev = l, l
	return l
}

func (l *EList) InsAfter(e *EList) *EList {
	e.next, e.prev = l.next, l
	l.next.prev, l.next = e, e
	return e
}

func (l *EList) InsBefore(e *EList) *EList {
	e.next, e.prev = l, l.prev
	l.prev.next, l.prev = e, e
	return e
}

func (l *EList) Next() *EList {
	return l.next
}

func (l *EList) Prev() *EList {
	return l.prev
}

func (l *EList) Ptr(fld uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(l)) - fld)
}

func (l *EList) Val() interface{} {
	return l
}
