package osam_simulator

import (
	"fmt"
	"log"
)

// ------------- Dummy PathORAM ------------- //
// see common.go for other type defs

type PathORAM struct {
	nl    int
	arr   [](map[int]Block)
	print bool
}

func CreateORAM(nleaves int, print bool) *PathORAM {
	me := &PathORAM{}
	me.nl = nleaves
	me.print = print
	me.arr = make([](map[int]Block), nleaves)
	for i := 0; i < nleaves; i++ {
		me.arr[i] = make(map[int]Block)
	}
	return me
}

func (oram *PathORAM) log(str string) {
	if oram.print {
		fmt.Println("[ORAM] " + str)
	}
}

// Access leaf (DUMMY)
func (oram *PathORAM) readRmAccess(a addr, callerMsg string) Block {
	i := a.leaf
	id := a.ctr
	if i >= oram.nl {
		log.Fatalf("ReadAndRm ACCESS leaf index out of bounds: i=%v, n=%v", i, oram.nl)
	}
	if callerMsg != "" {
		oram.log(fmt.Sprintf("ReadAndRm ACCESS: %v, called from: %v", a, callerMsg))
	} else {
		oram.log(fmt.Sprintf("ReadAndRm ACCESS: %v", a))
	}
	if v, ok := (oram.arr[i])[id]; ok {
		// need to "Remove" from the PathORAM leaf after reading
		delete(oram.arr[i], id)
		return v
	} else {
		oram.log(fmt.Sprintf("Read yielded None when reading %v", a))
		return Block{Data: NONE, IsNone: true}
	}
}

// NOTE: this is NOT the same functionality as [evict] in OSAM paper: we drop the stash for simulation;
// and we just directly write [value] into the leaf corresponding to [a].
// But this does not count as a separate "Access" of the ORAM: in real PathORAM implementation,
// [value] would just be placed on the LCA with the preceding read-Acess path address and [a].
func (oram *PathORAM) evictWrite(a addr, value interface{}) {
	oram.log(fmt.Sprintf("Evict=Write: storing value %v at %v", value, a))
	(oram.arr[a.leaf])[a.ctr] = Block{value, false}
}
