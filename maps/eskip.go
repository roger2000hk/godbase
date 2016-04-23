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

func (m *ESkip) Cut(start Iter, end Iter, fn TestFn) Sorted {
	if start == nil {
		start = m.root.next[ESkipLevels-1]
	}

	if end == nil {
		end = m.root.prev[ESkipLevels-1]
	}

	res := NewESkip()
	nn := &res.root

	for n := start.(*ESkipNode); n != end; n = n.next[ESkipLevels-1] {
		if n == &m.root {
			nn = &res.root
		} else if fn == nil || fn(n.key, n) {
			for i := 0; i < ESkipLevels-1; i++ {
				n.prev[i].next[i] = n.next[i]
				n.next[i].prev[i] = n.prev[i]			
				m.len--

				nn.next[i] = n
				n.prev[i] = nn
				n.next[i] = nn.next[i]
				nn.next[i].prev[i] = n
				nn = n
			}
			
			res.len++
		}		
	}

	return res
}

func (m *ESkip) Delete(key Cmp, val interface{}) int {
	cnt := 0

	if n, ok := m.FindNode(&m.root, key); ok {		
		for n.key == key {			
			next := n.next[ESkipLevels-1]
			
			if val == nil || n == val {
				n.Delete()
				cnt++
			}
			
			n = next
		}
	}

	m.len -= int64(cnt)
	return cnt
}

func (m *ESkip) First(start Iter, key Cmp) Iter {
	if start == nil {
		start = &m.root
	}

	if key == nil {
		return start
	}

	n, ok := m.FindNode(start.(*ESkipNode), key)

	if !ok {
		n = n.prev[ESkipLevels-1]
	}

	return n.prev[ESkipLevels-1]
}

func (m *ESkip) FindNode(start *ESkipNode, key Cmp) (*ESkipNode, bool) {
	next := start.next[ESkipLevels-1]
	if next != start && next != &m.root && key.Less(next.key) {
		return start, false
	}

	prev := m.root.prev[ESkipLevels-1]
	if prev != &m.root && prev.key.Less(key) {
		return prev, false
	}

	var pn *ESkipNode
	n := start
	maxSteps, steps := 1, 1

	for i := 0; i < ESkipLevels; i++ {
		n = n.next[i]

		for n != &m.root && n.key.Less(key){
			if steps == maxSteps && i > 0 {
				pn = pn.InsertAfter(n, i-1)
				steps = 0
			}

			n = n.next[i]
			steps++
		}

		if n.key == key {
			return n, true
		}

		n = n.prev[i]
		pn = n	
		steps = 1
		maxSteps++
	}

	return n, false
}

func (m *ESkip) Init() *ESkip {
	m.root.Init(nil)
	return m
}

func (m *ESkip) Insert(key Cmp, val interface{}, allowMulti bool) (Iter, bool) {
	n, ok := m.FindNode(&m.root, key)
	
	if ok && !allowMulti {
		return n, false
	}
	
	n.InsertAfter(val.(*ESkipNode), ESkipLevels-1)
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
		n.prev[i].next[i], n.next[i].prev[i] = n.next[i], n.prev[i] 
		n.prev[i], n.next[i] = n, n
	}
}

func (n *ESkipNode) HasNext() bool {
	return n.next[ESkipLevels-1].key != nil
}

func (n *ESkipNode) HasPrev() bool {
	return n.prev[ESkipLevels-1].key != nil
}

func (n *ESkipNode) Init(key Cmp) {
	n.key = key

	for i := 0; i < ESkipLevels; i++ {
		n.next[i], n.prev[i] = n, n
	}
}

func (n *ESkipNode) InsertAfter(node *ESkipNode, i int) *ESkipNode {
	node.prev[i], node.next[i] = n, n.next[i]
	n.next[i].prev[i], n.next[i] = node, node
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
