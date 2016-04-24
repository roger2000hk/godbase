package maps

import (
	"bytes"
	"fmt"
)

const ESkipLevels = 14

type ESkip struct {
	len int64
	root ESkipNode
}

func NewESkip() *ESkip {
	return new(ESkip).Init()
}

func (m *ESkip) Cut(start, end Iter, fn TestFn) Any {
	if start == nil {
		start = m.root.next[ESkipLevels-1]
	}

	if end == nil {
		end = &m.root
	}

	res := NewESkip()
	nn := &res.root

	for n := start.(*ESkipNode); n != end; n = n.next[ESkipLevels-1] {
		if Iter(n) == start {
			panic(fmt.Sprintf("invalid end: %v", end))
		}

		if n == &m.root {
			nn = &res.root
		} else if fn == nil || fn(n.key, n) {
			for i := 0; i < ESkipLevels-1; i++ {
				n.prev[i].next[i] = n.next[i]
				n.next[i].prev[i] = n.prev[i]			

				nn.next[i] = n
				n.prev[i] = nn

				n.next[i] = nn.next[i]
				nn.next[i].prev[i] = n
			}
			
			nn = n
			m.len--
			res.len++
		}		
	}

	return res
}

func (m *ESkip) Delete(start, end Iter, key Cmp, val interface{}) (Iter, int) {
	//fmt.Printf("delete1: %v\n", key)

	n := m.root.next[ESkipLevels-1]

	if start == nil {
		start = m.root.next[0]
	} else {
		n = start.(*ESkipNode)
	}

	if end == nil {
		end = &m.root
	}

	if key != nil {
		var ok bool
		if n, ok = m.FindNode(start, key); !ok {
			return n, 0
		}
	}

	cnt := 0
		
	for n != end && (key == nil || n == &m.root || n.key == key) {
		if n == start {
			panic(fmt.Sprintf("invalid end: %v", end))
		}

		next := n.next[ESkipLevels-1]
		
		if n != &m.root && (val == nil || n == val) {
			n.Delete()
			cnt++
		}
			
		n = next
	}

	m.len -= int64(cnt)
	return n.prev[ESkipLevels-1], cnt
}

func (m *ESkip) Find(start Iter, key Cmp, val interface{}) (Iter, bool) {
	n, ok := m.FindNode(start, key)

	if !ok {
		return n, false
	}
	
	for val != nil && n != val && n.key == key {
		n = n.next[ESkipLevels-1]
	}

	return n.prev[ESkipLevels-1], n.key == key && (val == nil || n == val)
}

func (m *ESkip) FindNode(start Iter, key Cmp) (*ESkipNode, bool) {
	if start == nil {
		start = m.root.next[0]
	}

	if next := m.root.next[ESkipLevels-1]; next != &m.root && key.Less(next.key) {
		return &m.root, false
	}

	if prev := m.root.prev[ESkipLevels-1]; prev != &m.root && prev.key.Less(key) {
		return prev, false
	}

	n := start.(*ESkipNode)
	var pn *ESkipNode
	maxSteps, steps := 1, 1

	for i := 0; i < ESkipLevels; i++ {
		for n != &m.root && n.key.Less(key)  {
			if steps == maxSteps && i > 0 {
				pn = pn.InsertAfter(n, i-1)
				steps = 0
			}
			
			n = n.next[i]
			steps++
		}
		
		if n.key == key {
			for n.prev[ESkipLevels-1].key == key {
				n = n.prev[ESkipLevels-1]
			}

			return n, true
		}
		
		n = n.prev[i]
		pn = n

		if i < ESkipLevels-1 {
			n = n.next[i+1]
		}

		steps = 1
		maxSteps++
	}

	return n, false
}

func (m *ESkip) Init() *ESkip {
	m.root.Init(nil)
	return m
}

func (m *ESkip) Insert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool) {
//	fmt.Printf("insert %v\n", key)

	n, ok := m.FindNode(start, key)
	
	if ok && !allowMulti {
		return n, false
	}
	n.InsertAfter(val.(*ESkipNode).Init(key), ESkipLevels-1)
 	m.len++
	return val.(*ESkipNode), true
}

func (m *ESkip) Len() int64 {
	return m.len
}

func (m *ESkip) String() string {
	var buf bytes.Buffer

	for i := 0; i < ESkipLevels; i++ {
		buf.WriteString("[")
		sep := ""

		for n := m.root.next[i]; n != &m.root; n = n.next[i] {
			fmt.Fprintf(&buf, "%v%v", sep, n.key)
			sep = ", "
		}

		buf.WriteString("]\n")
	}

	return buf.String()
}

type ESkipNode struct {
	key Cmp
	next [ESkipLevels]*ESkipNode
	prev [ESkipLevels]*ESkipNode
}

func (n *ESkipNode) Delete() {
	for i := 0; i < ESkipLevels; i++ {
		//fmt.Printf("delete2: %v / %v / %v\n", n.key, n.prev[i], n.next[i])

		n.prev[i].next[i] = n.next[i]
		n.next[i].prev[i] = n.prev[i] 
		n.prev[i], n.next[i] = nil, nil
	}
}

func (n *ESkipNode) HasNext() bool {
	return n.next[ESkipLevels-1].key != nil
}

func (n *ESkipNode) HasPrev() bool {
	return n.prev[ESkipLevels-1].key != nil
}

func (n *ESkipNode) Init(key Cmp) *ESkipNode {
	n.key = key

	for i := 0; i < ESkipLevels; i++ {
		n.next[i], n.prev[i] = n, n
	}

	return n
}

func (n *ESkipNode) InsertAfter(node *ESkipNode, i int) *ESkipNode {
	node.prev[i] = n
	node.next[i] = n.next[i]
	n.next[i].prev[i] = node
	n.next[i] = node
	return node
}

func (n *ESkipNode) Key() Cmp {
	return n.key
}

func (n *ESkipNode) Next() Iter {
	return n.next[ESkipLevels-1]
}

func (n *ESkipNode) Prev() Iter {
	return n.prev[ESkipLevels-1]
}

func (n *ESkipNode) Val() interface{} {
	return n
}
