package osam_simulator

import (
	"fmt"
	"log"
)

// ------------- Dummy ORAM ------------- //
// see common.go for other type defs

type PathORAM struct {
	nl    int
	arr   [](map[int]Block)
	debug bool
}

func CreateORAM(nleaves int, debug bool) *PathORAM {
	me := &PathORAM{}
	me.nl = nleaves
	me.debug = debug
	me.arr = make([](map[int]Block), nleaves)
	for i := 0; i < nleaves; i++ {
		me.arr[i] = make(map[int]Block)
	}
	return me
}

// Access leaf (DUMMY)
func (oram *PathORAM) Access(a addr, callerMsg string) Block {
	i := a.leaf
	id := a.ctr
	if i >= oram.nl {
		log.Fatalf("[ORAM] Access leaf index out of bounds: i=%v, n=%v", i, oram.nl)
	}
	if callerMsg != "" {
		fmt.Printf("[ORAM] Access: %v, called from: %v \n", i, callerMsg)
	} else {
		fmt.Printf("[ORAM] Access: %v \n", i)
	}

	if v, ok := (oram.arr[i])[id]; ok {
		delete(oram.arr[i], id)
		return v
	} else {
		return Block{Data: NONE, IsNone: true}
	}
}

//////// ORAM helpers //////////

func (oram *PathORAM) readAndRM(a addr, callerMsg string) Block {
	return oram.Access(a, callerMsg)
}

// NOTE: this is NOT the same functionality as evict in OSAM paper: bypass stash for simulation
func (oram *PathORAM) modEvict(a addr, value interface{}) {
	// note: this does NOT count as a new "Access" call:
	//  in real PathORAM, would just be placed on the LCA with the read-Access path address and a
	fmt.Printf("[ORAM] Evict: storing value %v at %v \n", value, a)
	(oram.arr[a.leaf])[a.ctr] = Block{value, false}
}
