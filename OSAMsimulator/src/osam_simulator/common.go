package osam_simulator

import (
	"strconv"
)

// ------------ BLOCK ------------
// INVARIANT for simulation: only use non-negative integers as values, so -1 = NONE
// type val = int
const NONE = -1

type Block struct {
	Data   interface{}
	IsNone bool
}

// ------------ ADDR ------------
// Based on OSAM paper: addresses include a global unique id "ctr" and the leaf index "leaf"
//
//	which is printed as "ctr_leaf"
type addr struct {
	ctr  int
	leaf int
}

var NIL = addr{NONE, NONE}

func (a addr) String() string {
	return strconv.Itoa(a.ctr) + "_" + strconv.Itoa(a.leaf)
}

// ------------ NODE (for SmartPointers) ------------
type Node struct {
	tailL   addr
	tailR   addr
	isRoot  bool
	content Block // invariant: if !isRoot, then content.isNone == true
	headP   addr  // invariant: if isRoot, then headP == NIL == {ctr:-1, leaf:-1}
}

// func (nd Node) String() string {
// 	return fmt.Sprintf("Node(%v)", nd)
// }
