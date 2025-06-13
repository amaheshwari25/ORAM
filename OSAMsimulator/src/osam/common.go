package osam

import "strconv"

const NONE = -1

// INVARIANT for simulation: only use non-negative integers as values, so -1 = NONE
type val = int

type block struct {
	data   interface{}
	isNone bool
}

// type block = interface{}

type addr struct {
	ctr  int
	leaf int
}

func (a addr) String() string {
	return strconv.Itoa(a.ctr) + "_" + strconv.Itoa(a.leaf)
}
