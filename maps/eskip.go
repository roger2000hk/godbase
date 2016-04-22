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

func (m *ESkip) Delete(key Cmp, val interface{}) int {
	cnt := 0

	if n, ok := m.FindNode(key); ok {		
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

func (m *ESkip) FindNode(key Cmp) (*ESkipNode, bool) {
	rootNext := m.root.next[ESkipLevels-1]

	if rootNext != &m.root {
		if key.Less(rootNext.key) {
			return &m.root, false
		}
		
		rootPrev := m.root.prev[ESkipLevels-1]
		if rootPrev.key.Less(key) {
			return rootPrev, false
		}
	}

	var pn *ESkipNode
	n := &m.root
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

func (m *ESkip) Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	n, ok := m.FindNode(key)
	
	if ok && !allowMulti {
		return n, false
	}
	
	n.InsertAfter(val.(*ESkipNode), ESkipLevels-1)
 	m.len++
	return val, true
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
