package storage

import (
	// "fmt"
	"sync"
)

type node struct {
	prev *node
	next *node
	key  string
	val  string
}

func NewNode(key, val string) *node {

	return &node{
		prev: nil,
		next: nil,
		key:  key,
		val:  val,
	}
}

type dll struct {
	sync.RWMutex
	head   *node
	tail   *node
	length int
	maxLgt int
}

func NewDll(maxLgt int) dll {
	return dll{
		head:   nil,
		tail:   nil,
		length: 0,
		maxLgt: maxLgt,
	}
}

func (d *dll) popNode() *node {

	// d.RLock()
	if d.length < 1 {
		return nil
	}
	// d.RUnlock()

	d.Lock()
	nd := d.tail
	if d.length > 1 {
		nTail := d.tail.prev
		d.tail.prev.next = nil
		d.tail = nTail
	} else if d.length == 1 {
		d.head = nil
		d.tail = nil
	}
	d.length--
	nd.prev = nil
	d.Unlock()

	return nd
}

func (d *dll) shiftNode() *node {

	// d.RLock()
	if d.length < 1 {
		return nil
	}
	// d.RUnlock()

	d.Lock()
	nd := d.head
	if d.length > 1 {
		nHead := d.head.next
		d.head.next.prev = nil
		d.head = nHead
	} else if d.length == 1 {
		d.head = nil
		d.tail = nil
	}
	d.length--
	nd.next = nil
	d.Unlock()

	return nd
}

// return: the first *node is the newNode, and the second one if not nil is the poped one
func (d *dll) unshift(key, val string) (*node, *node) {
	nn := NewNode(key, val)

	d.Lock()
	if d.length == 0 {
		d.head = nn
		d.tail = nn
		d.length++
	} else {
		// to avoid repetitive query to populate the dll
		if d.head.val != val {
			d.head.prev = nn
			nn.next = d.head
			d.head = nn
			d.length++
		}
	}
	d.Unlock()

	// d.RLock()
	if d.length > d.maxLgt {
		return nn, d.popNode()
	}
	// d.RUnlock()

	return nn, nil
}

// return first *node is the newNode second one, if not nil is the poped one
func (d *dll) unshiftNode(node *node) (*node, *node) {

	d.Lock()
	if d.length == 0 {
		d.head = node
		d.tail = node
	} else {
		d.head.prev = node
		node.next = d.head
		d.head = node
	}
	d.length++
	d.Unlock()

	// d.RLock()
	if d.length > d.maxLgt {
		return node, d.popNode()
	}
	// d.RUnlock()

	return node, nil
}

func (d *dll) removeNode(nd *node) *node {

	// d.RLock()
	if d.length == 1 || nd.next == nil {
		return d.popNode()
	}
	// d.RUnlock()

	// d.RLock()
	if nd.prev == nil {
		return d.shiftNode()
	}
	// d.RUnlock()

	d.Lock()
	nextNode := nd.next
	prevNode := nd.prev
	nextNode.prev = prevNode
	prevNode.next = nextNode

	nd.next = nil
	nd.prev = nil
	d.length--
	d.Unlock()

	return nd
}
