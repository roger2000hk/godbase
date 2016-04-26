package maps

import (
	"bytes"
	"fmt"
	"github.com/fncodr/godbase/lists"
	"unsafe"
)

type Skip struct {
	alloc *SkipAlloc
	bottom *SkipNode
	isInit bool
	len int64
	top SkipNode
}

func NewSkip(alloc *SkipAlloc, levels int) *Skip {
	return new(Skip).Init(alloc, levels)
}

func (m *Skip) AllocNode(key Key, val interface{}, prev *SkipNode) *SkipNode{
	if m.alloc == nil {
		 return new(SkipNode).Init(key, val, prev)
	} 
	
	return m.alloc.New(key, val, prev)
}

func (m *Skip) Cut(start, end Iter, fn MapFn) Any {
	if start == nil {
		start = m.bottom.next
	} else

	if end == nil {
		end = m.bottom.prev
	}

	res := NewSkip(m.alloc, m.Levels())
	nn := res.bottom

	n := start.(*SkipNode); 
	for Iter(n) != end {
		next := n.next
		
		if n == m.bottom {
			nn = res.bottom
		} else {
			k, v := n.key, n.val

			if fn != nil {
				k, v = fn(k, v)
			}

			if k != nil {
				for ln, lnn := n, nn;; ln = ln.up {
					ln.key, ln.val = k, v
					
					ln.prev.next = ln.next
					ln.next.prev = ln.prev
					
					ln.next = lnn.next
					lnn.next.prev = ln

					lnn.next = ln
					ln.prev = lnn					

					if ln.up == ln {
						break
					}
				}
				
				nn = n
				m.len--
				res.len++
			}
		}
		
		n = next
	}

	return res
}

func (m *Skip) Delete(start, end Iter, key Key, val interface{}) (Iter, int) {
	n := m.bottom.next

	if start == nil {
		start = m.top.next
	} else {
		n = start.(*SkipNode).Top()
	}

	if end == nil {
		end = m.bottom
	} else {
		end = end.(*SkipNode).Bottom()
	}

	if key != nil {
		var ok bool
		if n, ok = m.FindNode(start, key); !ok {
			return n, 0
		}
	} else {
		n = n.Bottom()
	}

	cnt := 0

	for n != end && (key == nil || n == m.bottom || n.key == key) {
		next := n.next

		if n != m.bottom && (val == nil || n.val == val) {
			n.Delete()
			cnt++

			if m.alloc != nil {
				m.alloc.Free(n)
			}
		}
		
		n = next
	}

	m.len -= int64(cnt)
	return n.prev, cnt
}

func (m *Skip) Find(start Iter, key Key, val interface{}) (Iter, bool) {
	n, ok := m.FindNode(start, key)
	
	if !ok {
		return n, false
	}

	for val != nil && n.key == key && n.val != val {
		n = n.next
	}
	
	return n.prev, n.key == key && (val == nil || n.val == val)
}

func (m *Skip) FindNode(start Iter, key Key) (*SkipNode, bool) {
	if start == nil {
		start = m.top.next
	}
	
	if next := m.bottom.next; next != m.bottom && key.Less(next.key) {
		return m.bottom, false
	}
		
	if prev := m.bottom.prev; prev != m.bottom && prev.key.Less(key) {
		return prev, false
	}

	var pn *SkipNode
	n := start.(*SkipNode)
	maxSteps, steps := 1, 1
	
	for true {
		if n.key == nil {
			n = n.next
		}

		for n.key != nil && n.key.Less(key) {
			if steps == maxSteps && pn != nil {
				var nn *SkipNode
				nn = m.AllocNode(n.key, n.val, pn)
				nn.down = n
				n.up, pn = nn, nn
				steps = 0
			}

			n = n.next
			steps++
		}

		if n.key == key {
			n = n.Bottom()

			for n.prev.key == key {
				n = n.prev
			}
			
			return n, true
		}

		pn = n.prev
		
		if pn.down == pn {
			break
		}

		n = pn.down

		steps = 1
		maxSteps++
	}
	
	return n.prev, false
}

func (m *Skip) Get(start Iter, key Key) (interface{}, bool) {
	n, ok := m.FindNode(start, key)
	
	if ok {
		return n.val, true
	}

	return nil, false
}

func (m *Skip) Init(alloc *SkipAlloc, levels int) *Skip {
	m.isInit = true
	m.alloc = alloc
	m.top.Init(nil, nil, nil)
	n := &m.top

	for i := 0; i < levels-1; i++ {
		n.down = m.AllocNode(nil, nil, nil)
		n.down.up = n
		n = n.down
	}

	n.down = n
	m.bottom = n
	return m
}

func (m *Skip) Insert(start Iter, key Key, val interface{}, allowMulti bool) (Iter, bool) {
	n, ok := m.FindNode(start, key)
	
	if ok && !allowMulti {
		n.val = val
		return n.prev, false
	}
	
	nn := m.AllocNode(key, val, n) 
	nn.down = nn
	m.len++
	return nn.prev, true
}

func (m *Skip) Len() int64 {
	return m.len
}

func (m *Skip) Levels() int {
	res := 1

	for n := &m.top; n.down != n; n = n.down { 
		res++
	}

	return res
}

func (m *Skip) String() string {
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
	freeNode lists.EDouble
	down, next, prev, up *SkipNode
	key Key
	val interface{}
}

var freeNodeOffs = unsafe.Offsetof(new(SkipNode).freeNode)

func (n *SkipNode) Bottom() *SkipNode {
	var res *SkipNode
	for res = n; res.down != res; res = res.down { }
	return res
}

func (n *SkipNode) Delete() {
	var pn *SkipNode

	for n != pn {
		n.prev.next, n.next.prev = n.next, n.prev
		pn = n
		n = n.up
	}
}

func (n *SkipNode) HasNext() bool {
	return n.next.key != nil
}

func (n *SkipNode) HasPrev() bool {
	return n.prev.key != nil
}

func (n *SkipNode) Init(key Key, val interface{}, prev *SkipNode) *SkipNode {
	n.key, n.val = key, val
	n.up = n

	if prev != nil {
		n.prev, n.next = prev, prev.next
		prev.next.prev, prev.next = n, n
	} else {
		n.prev, n.next = n, n
	}

	return n
}

func (n *SkipNode) Key() Key {
	return n.key
}

func (n *SkipNode) Next() Iter {
	return n.next
}

func (n *SkipNode) Prev() Iter {
	return n.prev
}

func (n *SkipNode) Top() *SkipNode {
	var res *SkipNode
	for res = n; res.up != res; res = res.up { }
	return res
}

func (n *SkipNode) Val() interface{} {
	return n.val
}

type SkipSlab []SkipNode

type SkipAlloc struct {
	freeList lists.EDouble
	idx int
	slab SkipSlab
	slabSize int
}

func NewSkipAlloc(slabSize int) *SkipAlloc {
	return new(SkipAlloc).Init(slabSize)
}

func (a *SkipAlloc) Init(slabSize int) *SkipAlloc {
	a.freeList.Init()
	a.slab = make(SkipSlab, slabSize)
	a.slabSize = slabSize
	return a
}

func (a *SkipAlloc) New(key Key, val interface{}, prev *SkipNode) *SkipNode {
	var res *SkipNode

	if n := a.freeList.Next(); n != n {
		n.Del()
		res = (*SkipNode)(n.Ptr(freeNodeOffs))
	} else {
		if a.idx == a.slabSize {
			a.slab = make([]SkipNode, a.slabSize)
			a.idx = 0
		}

		res = &a.slab[a.idx]
		a.idx++
	}

	return res.Init(key, val, prev)
}

func (a *SkipAlloc) Free(n *SkipNode) {
	a.freeList.InsAfter(&n.freeNode)
}
