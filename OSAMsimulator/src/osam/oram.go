package osam

import (
	"fmt"
	"log"
)

// ------------- Dummy ORAM ------------- //
// see common.go for other type defs

type PathORAM struct {
	nl    int
	arr   [](map[int]block)
	debug bool
}

func CreateORAM(nleaves int, debug bool) *PathORAM {
	me := &PathORAM{}
	me.nl = nleaves
	me.debug = debug
	me.arr = make([](map[int]block), nleaves)
	for i := 0; i < nleaves; i++ {
		me.arr[i] = make(map[int]block)
	}
	return me
}

// Access leaf (DUMMY)
func (oram *PathORAM) Access(a addr, callerMsg string) block {
	i := a.leaf
	id := a.ctr
	if i >= oram.nl {
		log.Fatalf("Access leaf index out of bounds: i=%v, n=%v", i, oram.nl)
	}
	if callerMsg != "" {
		fmt.Printf("ORAM Access: %v, called from: %v \n", i, callerMsg)
	} else {
		fmt.Printf("ORAM Access: %v \n", i)
	}

	if v, ok := (oram.arr[i])[id]; ok {
		return v
	} else {
		return block{data: NONE, isNone: true}
	}
}

//////// ORAM helpers //////////

func (oram *PathORAM) readAndRM(a addr, callerMsg string) block {
	return oram.Access(a, callerMsg)
}

// NOTE: this is NOT the same functionality as evict in OSAM paper: bypass stash for simulation
func (oram *PathORAM) modEvict(a addr, value interface{}) {
	// note: this does NOT count as a new "Access" call:
	//  in real PathORAM, would just be placed on the LCA with the read-Access path address and a
	fmt.Printf("Evict: storing value %v at %v \n", value, a)
	(oram.arr[a.leaf])[a.ctr] = block{value, false}
}
