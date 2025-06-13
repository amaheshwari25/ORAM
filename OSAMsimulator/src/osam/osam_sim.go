package osam

import (
	"fmt"
	"math/rand"
)

// --------- OSAM ------------ //
// see common.go for other type defs

type OSAM struct {
	counter int
	oram    *PathORAM
	// stash   []interface{}
}

type QueueNode struct {
	v    addr
	link addr
}

////////////////////////////////////////

func CreateOSAM(oram *PathORAM) *OSAM {
	o := &OSAM{}
	o.counter = 0
	o.oram = oram
	// o.stash = make([]block, stashSize)
	return o
}

func (osam *OSAM) Alloc() addr {
	leaf := rand.Intn(osam.oram.nl)
	a := addr{osam.counter, leaf}
	osam.counter++
	fmt.Printf("Alloc: %v \n", a)
	return a
}

func (osam *OSAM) Read(a addr) block {
	// 1. Read the value from address
	v := osam.oram.readAndRM(a, fmt.Sprintf("Read address %v", a))
	// 2. Don't actually need to do Evict in our dummy implementation
	// Evict
	return v
}

func GetData(b block) interface{} {
	return b.data
}

func (osam *OSAM) Write(a addr, value interface{}) {
	// 1. Simulate read Access by reading a dummy address
	osam.oram.readAndRM(osam.Alloc(), fmt.Sprintf("Write: %v @ address %v", value, a))
	// 2. Do the Evict (note: in dummy implementation, this directly places value v at addr a)
	osam.oram.modEvict(a, value)
}

///////////// QUEUE functionality ///////////////////

func (osam *OSAM) initQueue() (addr, addr) {
	head := osam.Alloc()
	tail := head
	return head, tail
}

func (osam *OSAM) enqueue(tail, a addr) addr {
	newTail := osam.Alloc()
	osam.Write(tail, QueueNode{v: a, link: newTail})
	return newTail

}

func (osam *OSAM) dequeue(head addr) (addr, addr) {
	b := osam.Read(head)
	if b.isNone {
		return addr{-1, -1}, addr{-1, -1}
	} else {
		bNode := b.data.(QueueNode) // FORCE data type to be a QueueNode
		return bNode.v, bNode.link
	}
}
