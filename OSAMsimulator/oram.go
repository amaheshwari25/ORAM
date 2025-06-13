package osamsimulator

import "fmt"

// ------------- Dummy ORAM ------------- //

type PathORAM struct {
	arr   []block
	debug bool
}

func CreateORAM(nleaves int, debug bool) *PathORAM {
	me := &PathORAM{}
	me.debug = debug
	me.arr = make([]block, nleaves)
	return me
}

// Access leaf (DUMMY)
func (oram *PathORAM) Access(i int) {
	fmt.Printf("ORAM Access: %v", i)
}
