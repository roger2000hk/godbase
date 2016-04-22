package godbase

import (
	"fmt"
	"bytes"
	"unsafe"
)

type SkipMap struct {
	bottom *SkipNode
	len int64
	nodeAlloc *SkipNodeAlloc
	top SkipNode
}

func NewSkipMap(nodeAlloc *SkipNodeAlloc, levels int) *SkipMap {
	return new(SkipMap).Init(nodeAlloc, levels)
}

func (m *SkipMap) Delete(key Cmp, val interface{}) int {
	cnt := 0

	if n, ok := m.FindNode(key); ok {
		for n.key == key {
			if val == nil || n.val == val {
				n.Delete()
				cnt++
				m.nodeAlloc.Free(n)
			}
			
			n = n.next
		}
	}

	m.len -= int64(cnt)
	return cnt
}

func (m *SkipMap) FindNode(key Cmp) (*SkipNode, bool) {
	if m.bottom.next != m.bottom {
		if key.Less(m.bottom.next.key) {
			return m.bottom, false
		}
		
		if m.bottom.prev.key.Less(key) {
			return m.bottom.prev, false
		}
	}

	var pn *SkipNode
	n := &m.top
	maxSteps, steps := 1, 1
	
	for true {
		n = n.next

		for n.key != nil && n.key.Less(key) {
			if steps == maxSteps && pn != nil {
				nn := m.nodeAlloc.New(n.key, n.val, pn)
				nn.down, n.up, pn = n, nn, nn
				steps = 0
			}

			n = n.next
			steps++
		}

		if n.key == key {
			return n, true
		}

		pn = n.prev

		if pn.down == pn {
			n = n.prev
			break
		}

		n = pn.down
		
		steps = 1
		maxSteps++
	}

	return n, false
}

func (m *SkipMap) Init(nodeAlloc *SkipNodeAlloc, levels int) *SkipMap {
	m.nodeAlloc = nodeAlloc
	m.top.Init()
	n := &m.top

	for i := 0; i < levels-1; i++ {
		n.down = nodeAlloc.New(nil, nil, nil)
		n.down.up = n
		n = n.down
	}

	n.down = n
	m.bottom = n
	return m
}

func (m *SkipMap) Insert(key Cmp, val interface{}, allowMulti bool) interface{} {
	n, ok := m.FindNode(key)
	
	if ok && !allowMulti {
		return n.val
	}
	
	nn := m.nodeAlloc.New(key, val, n) 
	nn.down = nn
	m.len++
	return val
}

func (m *SkipMap) Len() int64 {
	return m.len
}

func (m *SkipMap) String() string {
	var buf bytes.Buffer
	start := &m.top

	for true {
		buf.WriteString("[")
		sep := ""

		for n := start.next; n.key != nil; n = n.next {
			fmt.Fprintf(&buf, "%v%v", sep, n.key)
			if n.val != nil {
				fmt.Fprintf(&buf, ": %v", n.val)
			}
			sep = ", "
		}

		buf.WriteString("]\n")

		if start.down == start {
			break
		}

		start = start.down
	}

	return buf.String()
}

type SkipNode struct {
	freeNode EList
	down, next, prev, up *SkipNode
	key Cmp
	val interface{}
}

var freeNodeOffs = unsafe.Offsetof(new(SkipNode).freeNode)

func (n *SkipNode) Delete() {
	var pn *SkipNode

	for n != pn {
		n.prev.next, n.next.prev = n.next, n.prev
		pn = n
		n = n.up
	}
}

func (n *SkipNode) Init() {
	n.down, n.next, n.prev, n.up = n, n, n, n
}

type SkipSlab []SkipNode

type SkipNodeAlloc struct {
	freeList EList
	idx int
	slab SkipSlab
	slabSize int
}

func NewSkipNodeAlloc(slabSize int) *SkipNodeAlloc {
	return new(SkipNodeAlloc).Init(slabSize)
}

func (a *SkipNodeAlloc) Init(slabSize int) *SkipNodeAlloc {
	a.freeList.Init()
	a.slab = make(SkipSlab, slabSize)
	a.slabSize = slabSize
	return a
}

func (a *SkipNodeAlloc) New(key Cmp, val interface{}, prev *SkipNode) *SkipNode {
	var res *SkipNode

	if n := a.freeList.next; n != n {
		n.Del()
		res = (*SkipNode)(n.Ptr(freeNodeOffs))
	}

	if res == nil {
		if a.idx == a.slabSize {
			a.slab = make([]SkipNode, a.slabSize)
			a.idx = 0
		}

		res = &a.slab[a.idx]
		a.idx++
	}

	res.key, res.val = key, val
	res.up = res

	if prev != nil {
		res.prev, res.next = prev, prev.next
		prev.next.prev, prev.next = res, res
	} else {
		res.prev, res.next = res, res
	}

	return res
}

func (a *SkipNodeAlloc) Free(n *SkipNode) {
	a.freeList.InsAfter(&n.freeNode)
}
