package maps

import (
	"bytes"
	"fmt"
	"github.com/fncodr/godbase"
)

const ELevels = 10

type ESort struct {
	isInit bool
	len int64
	root ENode
}

type ENode struct {
	key godbase.Key
	next [ELevels]*ENode
	prev [ELevels]*ENode
}

func NewESort() *ESort {
	return new(ESort).Init()
}

func (m *ESort) Clear() {
	for i := 0; i < ELevels; i++ {
		m.root.next[i], m.root.prev[i] = nil, nil
	}

	m.len = 0
}

func (m *ESort) Cut(start, end godbase.Iter, fn godbase.KVMapFn) godbase.Map {
	if start == nil {
		start = m.root.next[ELevels-1]
	}

	if end == nil {
		end = &m.root
	}

	res := NewESort()
	n, nn := start.(*ENode), &res.root

	for n != end  {
		next := n.next[ELevels-1]

		if n == &m.root {
			nn = &res.root
		} else {
			k, v := n.key, interface{}(n)
			
			if fn != nil {
				k, v = fn(k, v)
			}

			if k != nil {
				vn := v.(*ENode)
				vn.key = k

				for i := 0; i < ELevels; i++ {
					if n.next[i] != n {
						n.prev[i].next[i] = n.next[i]
						n.next[i].prev[i] = n.prev[i]			

						vn.next[i] = nn.next[i]
						nn.next[i].prev[i] = vn
						nn.next[i] = vn
						vn.prev[i] = nn	
					}
				}
				
				nn = vn
				m.len--
				res.len++
			}
		}

		n = next
	}

	return res
}

func (m *ESort) Delete(start, end godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, int) {
	n := m.root.next[ELevels-1]

	if start == nil {
		start = m.root.next[0]
	} else {
		n = start.(*ENode)
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
		next := n.next[ELevels-1]
		
		if n != &m.root && (val == nil || n == val) {
			n.Delete()
			cnt++
		}
			
		n = next
	}

	m.len -= int64(cnt)
	return n, cnt
}

func (m *ESort) Find(start godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, bool) {
	n, ok := m.FindNode(start, key)

	if !ok {
		return n, false
	}
	
	for val != nil && n != val && n.key == key {
		n = n.next[ELevels-1]
	}

	return n, n.key == key && (val == nil || n == val)
}

func (m *ESort) FindNode(start godbase.Iter, key godbase.Key) (*ENode, bool) {
	if start == nil {
		start = m.root.next[0]
	}

	if next := m.root.next[ELevels-1]; next != &m.root && key.Less(next.key) {
		return &m.root, false
	}

	if prev := m.root.prev[ELevels-1]; prev != &m.root && prev.key.Less(key) {
		return prev, false
	}

	n := start.(*ENode)
	var pn *ENode
	maxSteps, steps := 1, 1

	for i := 0; i < ELevels; i++ {
		isless := false

		if n != &m.root {
			isless = n.key.Less(key)
		}
		
		for isless  {
			if steps == maxSteps && i > 0 {
				pn = pn.InsertAfter(n, i-1)
				steps = 0
			}
			
			n = n.next[i]
			isless = n != &m.root && n.key.Less(key)
			steps++
		}
		
		if !isless && n.key == key {
			for n.prev[ELevels-1].key == key {
				n = n.prev[ELevels-1]
			}

			return n, true
		}
		
		n = n.prev[i]
		pn = n

		if i < ELevels-1 {
			n = n.next[i+1]
		}

		steps = 1
		maxSteps++
	}

	return n, false
}

func (m *ESort) First() godbase.Iter {
	return m.root.next[ELevels-1]
}

func (m *ESort) Get(key godbase.Key) (interface{}, bool) {
	n, ok := m.FindNode(nil, key)
	
	if ok {
		return n, true
	}

	return nil, false
}

func (m *ESort) Init() *ESort {
	m.isInit = true
	m.root.Init(nil)
	return m
}

func (m *ESort) Insert(start godbase.Iter, key godbase.Key, val interface{}, 
	allowMulti bool) (godbase.Iter, bool) {
	n, ok := m.FindNode(start, key)
	
	if ok && !allowMulti {
		return n, false
	}
	n.InsertAfter(val.(*ENode).Init(key), ELevels-1)
 	m.len++
	return val.(*ENode), true
}

func (m *ESort) Len() int64 {
	return m.len
}

func (m *ESort) New() godbase.Map {
	return NewESort()
}

func (m *ESort) Set(key godbase.Key, val interface{}) bool {
	i, ok := m.Insert(nil, key, val, false)

	if !ok && i != val {
		panic("changing embedded val for key is not supported!")
	}

	return ok
}

func (m *ESort) String() string {
	var buf bytes.Buffer

	for i := 0; i < ELevels; i++ {
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

func (n *ENode) Delete() {
	for i := 0; i < ELevels; i++ {
		n.prev[i].next[i] = n.next[i]
		n.next[i].prev[i] = n.prev[i] 
		n.prev[i], n.next[i] = nil, nil
	}
}

func (n *ENode) Init(key godbase.Key) *ENode {
	n.key = key

	for i := 0; i < ELevels; i++ {
		n.next[i], n.prev[i] = n, n
	}

	return n
}

func (n *ENode) InsertAfter(node *ENode, i int) *ENode {
	node.prev[i] = n
	node.next[i] = n.next[i]
	n.next[i].prev[i] = node
	n.next[i] = node
	return node
}

func (n *ENode) Key() godbase.Key {
	return n.key
}

func (n *ENode) Next() godbase.Iter {
	return n.next[ELevels-1]
}

func (n *ENode) Val() interface{} {
	return n
}

func (n *ENode) Valid() bool {
	return n.key != nil
}

func (m *ESort) While(fn godbase.KVTestFn) bool {
	for n := m.root.next[ELevels-1]; n != &m.root; n = n.next[ELevels-1] {
		if !fn(n.key, n) {
			return false
		}
	} 

	return true
}
