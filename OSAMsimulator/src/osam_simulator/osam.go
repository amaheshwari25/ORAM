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
	print   bool
	reads   map[addr]bool
	writes  map[addr]bool
	allocs  map[addr]bool
}

func (osam *OSAM) log(str string) {
	if osam.print {
		fmt.Println("[OSAM] " + str)
	}
}

type QueueElem struct {
	v    addr
	link addr
}

func (qe QueueElem) String() string {
	return fmt.Sprintf("QE(v:%v->l:%v)", qe.v, qe.link)
}

////////////////////////////////////////

func CreateOSAM(oram *PathORAM, print bool) *OSAM {
	o := &OSAM{}
	o.counter = 0
	o.oram = oram
	o.print = print
	o.reads = make(map[addr]bool)
	o.writes = make(map[addr]bool)
	o.allocs = make(map[addr]bool)
	return o
}

func (osam *OSAM) Alloc(msg string) addr {
	leaf := rand.Intn(osam.oram.nl)
	a := addr{osam.counter, leaf}
	osam.counter++
	osam.allocs[a] = true
	osam.log(fmt.Sprintf("Alloc: %v for %v", a, msg))
	return a
}

func (osam *OSAM) Read(a addr) Block {
	assert(hasAddr(osam.allocs, a), fmt.Sprintf("Address %v has not been alloc'd", a))
	assert(!hasAddr(osam.reads, a), fmt.Sprintf("Address %v has already been read", a))
	osam.reads[a] = true
	// 1. Read the value from address
	v := osam.oram.readRMAccess(a, fmt.Sprintf("Read address %v", a))
	// 2. Don't actually need to do Evict in our dummy implementation
	return v
}

func (osam *OSAM) Write(a addr, value interface{}, msg string) {
	assert(hasAddr(osam.allocs, a), fmt.Sprintf("Address %v has not been alloc'd", a))
	assert(!hasAddr(osam.writes, a), fmt.Sprintf("Address %v has already been written to", a))
	osam.writes[a] = true
	// 1. Simulate Read Access by reading a dummy address
	osam.oram.readRMAccess(osam.Alloc(fmt.Sprintf("Write at addr %v (DUMMY)", a)), msg)
	// 2. Do the Evict (in this dummy implementation, this directly places value at addr a)
	osam.oram.modEvict(a, value)
}

func (osam *OSAM) writeQE(a addr, value QueueElem) {
	msg := fmt.Sprintf("Write: %v @ address %v", value, a)
	osam.Write(a, value, msg)
}

func (osam *OSAM) writeN(a addr, value *Node) {
	msg := fmt.Sprintf("Write: %v @ address %v", *value, a)
	osam.Write(a, value, msg)
}

///////////// QUEUE functionality ///////////////////

func (osam *OSAM) initQueue() (addr, addr) {
	head := osam.Alloc("initQueue")
	tail := head
	return head, tail
}

func (osam *OSAM) enqueue(tail, a addr) addr {
	newTail := osam.Alloc(fmt.Sprintf("enqueue addr %v at tail %v", a, tail))
	osam.writeQE(tail, QueueElem{v: a, link: newTail})
	return newTail

}

func (osam *OSAM) dequeue(head addr) (addr, addr) {
	b := osam.Read(head)
	if b.IsNone {
		return NIL, NIL
	} else {
		bNode := b.Data.(QueueElem) // force data type to be a QueueElem
		return bNode.v, bNode.link
	}
}
