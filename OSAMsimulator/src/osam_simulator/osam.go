package osam_simulator

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

type QueueElem struct {
	v    addr
	link addr
}

func (qe QueueElem) String() string {
	return fmt.Sprintf("QE(v:%v->l:%v)", qe.v, qe.link)
}

////////////////////////////////////////

func CreateOSAM(oram *PathORAM) *OSAM {
	o := &OSAM{}
	o.counter = 0
	o.oram = oram
	// o.stash = make([]Block, stashSize)
	return o
}

func (osam *OSAM) Alloc(msg string) addr {
	leaf := rand.Intn(osam.oram.nl)
	a := addr{osam.counter, leaf}
	osam.counter++
	fmt.Printf("[OSAM] Alloc: %v for %v \n", a, msg)
	return a
}

func (osam *OSAM) Read(a addr) Block {
	// 1. Read the value from address
	v := osam.oram.readAndRM(a, fmt.Sprintf("Read address %v", a))
	// 2. Don't actually need to do Evict in our dummy implementation
	// Evict
	return v
}

func (osam *OSAM) WriteQE(a addr, value QueueElem) {
	msg := fmt.Sprintf("Write: %v @ address %v", value, a)
	osam.Write(a, value, msg)
}

func (osam *OSAM) WriteN(a addr, value *Node) {
	msg := fmt.Sprintf("Write: %v @ address %v", *value, a)
	osam.Write(a, value, msg)
}

func (osam *OSAM) Write(a addr, value interface{}, msg string) {
	// 1. Simulate read Access by reading a dummy address
	osam.oram.readAndRM(osam.Alloc(fmt.Sprintf("Write at addr %v", a)), msg)
	// 2. Do the Evict (note: in dummy implementation, this directly places value v at addr a)
	osam.oram.modEvict(a, value)
}

///////////// QUEUE functionality ///////////////////

func (osam *OSAM) initQueue() (addr, addr) {
	head := osam.Alloc("initQueue")
	tail := head
	return head, tail
}

func (osam *OSAM) enqueue(tail, a addr) addr {
	newTail := osam.Alloc(fmt.Sprintf("enqueue addr %v at tail %v", a, tail))
	osam.WriteQE(tail, QueueElem{v: a, link: newTail})
	return newTail

}

func (osam *OSAM) dequeue(head addr) (addr, addr) {
	b := osam.Read(head)
	if b.IsNone {
		return NIL, NIL
	} else {
		bNode := b.Data.(QueueElem) // FORCE data type to be a QueueElem
		return bNode.v, bNode.link
	}
}
