package maps

import (
	"bytes"
	"fmt"
	"github.com/fncodr/godbase/lists"
	"unsafe"
)

type Sort struct {
	alloc *SlabAlloc
	bottom *Node
	isInit bool
	len int64
	top Node
}

type SlabAlloc struct {
	freeList lists.EDouble
	idx int
	slab Slab
	slabSize int
}

type Slab []Node

func NewSlab(a *SlabAlloc, ls int) *Sort {
	return new(Sort).Init(a, ls)
}

func NewSort(ls int) *Sort {
	return NewSlab(nil, ls)
}

func (m *Sort) SlabAllocNode(key Key, val interface{}, prev *Node) *Node{
	if m.alloc == nil {
		 return new(Node).Init(key, val, prev)
	} 
	
	return m.alloc.New(key, val, prev)
}

func (n *Node) Bottom() *Node {
	var res *Node
	for res = n; res.down != res; res = res.down { }
	return res
}

func (m *Sort) Clear() {
	if m.alloc != nil {
		for n := m.bottom.next; n != m.bottom; n = n.next {
			m.alloc.Free(n)
		}
	}

	for n := &m.top;; n = n.next {
		n.next, n.prev = n, n

		if n == m.bottom {
			break
		}
	}

	m.len = 0
}

func (m *Sort) Cut(start, end Iter, fn MapFn) Any {
	if start == nil {
		start = m.bottom.next
	} else

	if end == nil {
		end = m.bottom.prev
	}

	res := NewSlab(m.alloc, m.Levels())
	nn := res.bottom

	n := start.(*Node); 
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

func (m *Sort) Delete(start, end Iter, key Key, val interface{}) (Iter, int) {
	n := m.bottom.next

	if start == nil {
		start = m.top.next
	} else {
		n = start.(*Node).Top()
	}

	if end == nil {
		end = m.bottom
	} else {
		end = end.(*Node).Bottom()
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
	return n, cnt
}

func (n *Node) Delete() {
	var pn *Node

	for n != pn {
		n.prev.next, n.next.prev = n.next, n.prev
		pn = n
		n = n.up
	}
}

func (m *Sort) Find(start Iter, key Key, val interface{}) (Iter, bool) {
	n, ok := m.FindNode(start, key)
	
	if !ok {
		return n, false
	}

	for val != nil && n.key == key && n.val != val {
		n = n.next
	}
	
	return n, n.key == key && (val == nil || n.val == val)
}

func (m *Sort) FindNode(start Iter, key Key) (*Node, bool) {
	if start == nil {
		start = m.top.next
	}
	
	if next := m.bottom.next; next != m.bottom && key.Less(next.key) {
		return m.bottom, false
	}
		
	if prev := m.bottom.prev; prev != m.bottom && prev.key.Less(key) {
		return prev, false
	}

	var pn *Node
	n := start.(*Node)
	maxSteps, steps := 1, 1
	
	for true {
		if n.key == nil {
			n = n.next
		}

		isless := false
		if n.key != nil {
			isless = n.key.Less(key)
		}

		for isless {
			if steps == maxSteps && pn != nil {
				var nn *Node
				nn = m.SlabAllocNode(n.key, n.val, pn)
				nn.down = n
				n.up, pn = nn, nn
				steps = 0
			}

			n = n.next
			isless = n.key != nil && n.key.Less(key)
			steps++
		}

		if !isless && n.key == key {
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

func (m *Sort) First() Iter {
	return m.bottom.next
}


func (m *Sort) Get(key Key) (interface{}, bool) {
	n, ok := m.FindNode(nil, key)
	
	if ok {
		return n.val, true
	}

	return nil, false
}

func (a *SlabAlloc) Free(n *Node) {
	a.freeList.InsAfter(&n.freeNode)
}

func (m *Sort) Init(alloc *SlabAlloc, levels int) *Sort {
	m.isInit = true
	m.alloc = alloc
	m.top.Init(nil, nil, nil)
	n := &m.top

	for i := 0; i < levels-1; i++ {
		n.down = m.SlabAllocNode(nil, nil, nil)
		n.down.up = n
		n = n.down
	}

	n.down = n
	m.bottom = n
	return m
}

func (m *Sort) Insert(start Iter, key Key, val interface{}, allowMulti bool) (Iter, bool) {
	n, ok := m.FindNode(start, key)
	
	if ok && !allowMulti {
		return n, false
	}
	
	nn := m.SlabAllocNode(key, val, n) 
	nn.down = nn
	m.len++

	return nn, true
}

func (m *Sort) Len() int64 {
	return m.len
}

func (m *Sort) Levels() int {
	res := 1

	for n := &m.top; n.down != n; n = n.down { 
		res++
	}

	return res
}

func (m *Sort) New() Any {
	return NewSlab(m.alloc, m.Levels())
}

func (m *Sort) Set(key Key, val interface{}) bool {
	i, ok := m.Insert(nil, key, val, false)

	if !ok {
		i.(*Node).val = val
	}

	return ok
}

func (m *Sort) String() string {
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

type Node struct {
	freeNode lists.EDouble
	down, next, prev, up *Node
	key Key
	val interface{}
}

var freeNodeOffs = unsafe.Offsetof(new(Node).freeNode)

func (n *Node) Init(key Key, val interface{}, prev *Node) *Node {
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

func (n *Node) Key() Key {
	return n.key
}

func (n *Node) Next() Iter {
	return n.next
}

func (n *Node) Prev() Iter {
	return n.prev
}

func (n *Node) Top() *Node {
	var res *Node
	for res = n; res.up != res; res = res.up { }
	return res
}

func (n *Node) Val() interface{} {
	return n.val
}

func (n *Node) Valid() bool {
	return n.key != nil
}

func NewSlabAlloc(slabSize int) *SlabAlloc {
	return new(SlabAlloc).Init(slabSize)
}

func (a *SlabAlloc) Init(slabSize int) *SlabAlloc {
	a.freeList.Init()
	a.slab = make(Slab, slabSize)
	a.slabSize = slabSize
	return a
}

func (a *SlabAlloc) New(key Key, val interface{}, prev *Node) *Node {
	var res *Node

	if n := a.freeList.Next(); n != n {
		n.Del()
		res = (*Node)(n.Ptr(freeNodeOffs))
	} else {
		if a.idx == a.slabSize {
			a.slab = make(Slab, a.slabSize)
			a.idx = 0
		}

		res = &a.slab[a.idx]
		a.idx++
	}

	return res.Init(key, val, prev)
}

func (m *Sort) While(fn TestFn) bool {
	for n := m.bottom.next; n != m.bottom; n = n.next {
		if !fn(n.key, n.val) {
			return false
		}
	} 

	return true
}
