package osam_simulator

import (
	"strconv"
)

func assert(cond bool, msg string) {
	if !cond {
		panic("Assert failed: " + msg)
	}
}

func hasAddr(m map[addr]bool, a addr) bool {
	_, ok := m[a]
	return ok
}

// for suppressing / un-supressing all print output
var suppressPrint = false

func Suppress() {
	suppressPrint = true
}

func Unsupress() {
	suppressPrint = false
}

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
// which is printed as "ctr_leaf"
type addr struct {
	ctr  int
	leaf int
}

var NIL = addr{NONE, NONE}

func (a addr) String() string {
	return strconv.Itoa(a.ctr) + "_" + strconv.Itoa(a.leaf)
}

// ------------ Ptr (for all SmartPointer implementations) ------------
type Ptr struct {
	head addr
}

// ------------ Node (for base SmartPointers) ------------
type Node struct {
	tailL   addr
	tailR   addr
	isRoot  bool
	content Block // invariant: if !isRoot, then content.isNone == true
	headP   addr  // invariant: if isRoot, then headP == NIL == {ctr:-1, leaf:-1}

	// My new metadata (not in paper)
	id int
}

// ------------ BNode (for Balanced SmartPointers) ------------
type BNode struct {
	tailL   addr
	headL   addr
	tailR   addr
	headR   addr
	isRoot  bool
	content Block // invariant: if !isRoot, then content.isNone == true
	count   int   // invariant: if !isRoot, then count == NONE == -1
	headP   addr  // invariant: if isRoot, then headP == NIL == {ctr:-1, leaf:-1}
	tailP   addr  // invariant: if isRoot, then tailP == NIL == {ctr:-1, leaf:-1}

	// My new metadata (not in paper)
	id int
}
