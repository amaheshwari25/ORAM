package osam_simulator

import (
	"fmt"
	"log"
)

// ------------- Dummy ORAM ------------- //
// see common.go for other type defs

type PathORAM struct {
	nl        int
	arr       [](map[int]Block)
	printORAM bool
}

func CreateORAM(nleaves int, print bool) *PathORAM {
	me := &PathORAM{}
	me.nl = nleaves
	me.printORAM = print
	me.arr = make([](map[int]Block), nleaves)
	for i := 0; i < nleaves; i++ {
		me.arr[i] = make(map[int]Block)
	}
	return me
}

// Access leaf (DUMMY)
func (oram *PathORAM) readRMAccess(a addr, callerMsg string) Block {
	i := a.leaf
	id := a.ctr
	if i >= oram.nl {
		log.Fatalf("[ORAM] ReadAndRM Access leaf index out of bounds: i=%v, n=%v", i, oram.nl)
	}
	if oram.printORAM {
		if callerMsg != "" {
			fmt.Printf("[ORAM] ReadAndRM Access: %v, called from: %v \n", i, callerMsg)
		} else {
			fmt.Printf("[ORAM] ReadAndRM Access: %v \n", i)
		}
	}

	if v, ok := (oram.arr[i])[id]; ok {
		delete(oram.arr[i], id)
		return v
	} else {
		fmt.Printf("[ORAM] Read NONE when reading %v \n", i)
		return Block{Data: NONE, IsNone: true}
	}
}

//////// ORAM helpers //////////

// func (oram *PathORAM) readAndRM(a addr, callerMsg string) Block {
// 	return oram.Access(a, callerMsg)
// }

// NOTE: this is NOT the same functionality as evict in OSAM paper: bypass stash for simulation
func (oram *PathORAM) modEvict(a addr, value interface{}) {
	// note: this does NOT count as a new "Access" call:
	//  in real PathORAM, would just be placed on the LCA with the read-Access path address and a
	if oram.printORAM {
		fmt.Printf("[ORAM] Evict: storing value %v at %v \n", value, a)
	}
	(oram.arr[a.leaf])[a.ctr] = Block{value, false}
}
