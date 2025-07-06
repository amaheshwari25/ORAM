package osam_simulator

// Algorithm to construct the emulated graph in OSAM paper from a adjacency-list graph representation

import (
	"fmt"
	"math"
	"sort"
)

// ** Current notes ** :
//  - not actually using OSAM right now, dummy implementation for PoC
//  - no actual dummy operations happening

// -------- TYPE DEFINITIONS --------- //

type ptr = int

type VtxType = int

const (
	Real VtxType = iota
	Inc
	Out
	Internal
)

type InputEdge struct {
	U int
	V int
	W int // all real weights guaranteed to be positive
}

type InputGraph struct {
	print bool
	edges []InputEdge
}

type Vtx struct {
	Id    int // Real vertex that this is associated with
	Type  VtxType
	Other int // If VtxType = Inc or Out: Other = other part of this edge
	UP    ptr
	LC    ptr
	RC    ptr
}

type Edge struct {
	Ptr    ptr
	Weight int
}

type OSAMGraph struct {
	print   bool
	pCtr    ptr // global counter for assigning dummy "pointers" = ptr
	in_deg  []int
	out_deg []int

	Vtcs    []ptr
	FakeRAM map[ptr](*Vtx)
}

// -------- CONSTRUCTOR FUNCTIONS --------- //
func CreateInpEdge(u, v, w int) InputEdge {
	return InputEdge{u, v, w}
}

func CreateInputGraph(print bool, us []int, vs []int, ws []int) *InputGraph {
	assert(len(us) == len(vs), "mismatched input lengths")
	assert(len(us) == len(ws), "mismatched input lengths")
	edges := make([]InputEdge, len(us))
	for i := 0; i < len(us); i++ {
		edges[i] = CreateInpEdge(us[i], vs[i], ws[i])
	}
	return &InputGraph{print, edges}
}

func (inpG *InputGraph) CreateOSAMGraph() *OSAMGraph {
	return inpG.construct()
}

// -------- HELPER FUNCTIONS --------- //

func (inpG *InputGraph) log(str string, newline bool) {
	if newline {
		fmt.Println()
	}
	if inpG.print {
		fmt.Println("[IG] " + str)
	}
}

// Returns the smallest power of 2 >= x
func nextPowTwo(x int) int {
	p := 1
	for p < x {
		p <<= 1 // left shift until > 2
	}
	return p
}

// Returns the max j s.t. 2^j divides x
func maxDivPowTwo(x int) int {
	p := 1
	ctr := -1
	for p <= x && x%p == 0 {
		p <<= 1
		ctr++
	}
	return ctr
}

// Returns ceil(lg(x))
func lg(x int) int {
	return int(math.Ceil(math.Log2(float64(x))))
}

// Dummy sort comparisons for an O-sort procedure.
// This will just sorts normally; assume replaced by O-sort.
// TBD fix: bad design right now with sort.Slice
func (inpG *InputGraph) compareU(i, j int) bool {
	cmpU := inpG.edges[i].U - inpG.edges[j].U
	if cmpU == 0 {
		cmpW := inpG.edges[i].W - inpG.edges[j].W
		if cmpW == 0 {
			return inpG.edges[i].V < inpG.edges[j].V
		}
		return cmpW < 0
	}
	return cmpU < 0
}

func (inpG *InputGraph) compareV(i, j int) bool {
	cmpV := inpG.edges[i].V - inpG.edges[j].V
	if cmpV == 0 {
		cmpW := inpG.edges[i].W - inpG.edges[j].W
		if cmpW == 0 {
			return inpG.edges[i].U < inpG.edges[j].U
		}
		return cmpW < 0
	}
	return cmpV < 0
}

// TBD: for full implementation, add compareW for final O-Sort

func (inpG *InputGraph) computeDegs(deg *[]int) {
	l := len(inpG.edges)
	arr := *deg
	assert(len(arr) == l, "mismatched lengths")
	cumu := 0
	for i := l - 1; i >= 0; i-- {
		if inpG.edges[i].W == NONE { // Vertex
			arr[i] = cumu
			cumu = 0
		} else {
			arr[i] = NONE
			cumu++
		}
	}
}

// Return the address [ptr] at which vertex was created
func (oG *OSAMGraph) createVtx(id int, vtxType VtxType, other int, up ptr, lc ptr, rc ptr) int {
	oG.pCtr++
	v := Vtx{Id: id, Type: vtxType, Other: other, UP: up, LC: lc, RC: rc}
	oG.FakeRAM[oG.pCtr] = &v
	return oG.pCtr
}

////////////////////////////////////////

func (inpG *InputGraph) construct() *OSAMGraph {
	l := len(inpG.edges)

	osamG := OSAMGraph{
		print: inpG.print, pCtr: 0,
		in_deg: make([]int, l), out_deg: make([]int, l),
		FakeRAM: make(map[ptr]*Vtx)}
	// TBD: set the Vtcs variable

	// 1. O-SORT: inpG edges by v (head vertex)
	sort.Slice(inpG.edges, inpG.compareV)
	inpG.log(fmt.Sprintf("%v", inpG.edges), true)
	// 2. LINEAR-SCAN: Compute in_deg array
	inpG.computeDegs(&osamG.in_deg)
	inpG.log(fmt.Sprintf("%v", osamG.in_deg), true)
	// 3. LINEAR-SCAN + binary-pointer-tree: create inc_vtcs array
	osamG.Vtcs = inpG.createIncTrees(&osamG)

	// // 4. O-SORT: edges by u, inc_vtcs by u
	// sort.Slice(inpG.edges, inpG.compareV)
	// sort.Slice() // TBD: will need to fix sorting

	// // 5. LINEAR-SCAN: Compute out_deg array
	// inpG.computeDegs(&osamG.out_deg)
	// inpG.log(fmt.Sprintf("Out-degrees: %v", osamG.out_deg), true)

	// 6. LINEAR-SCAN + binary-pointer-tree: compute out_vtcs array

	return &osamG
}

// Binary-pointer-tree construction in a streaming pass over the leaves.
// Requires O(log E) client storage, indicated by the client array C (and O(1)-size variables)

func (inpG *InputGraph) createIncTrees(osamG *OSAMGraph) []ptr {
	l := len(inpG.edges)
	assert(l == len(osamG.in_deg), "mismatched lengths")
	assert(osamG.in_deg[0] != NONE, "in_deg not created properly")
	store := osamG.FakeRAM

	inc_vtcs := make([]ptr, l)

	// Client state
	m := NONE
	z := NONE
	v := NONE
	startInd := NONE
	C := make([]int, lg(l)+1)
	for k := 0; k < len(C); k++ {
		C[k] = NONE
	}

	// Streaming pass
	for i := 0; i < l; i++ {
		// fmt.Printf("Edge list ind: %v \n", i)
		id := inpG.edges[i].V
		if osamG.in_deg[i] != NONE { // Vertex encountered: start new tree creation
			// fmt.Printf("Starting: state of C: %v \n", C)
			m = osamG.in_deg[i]
			z = 2*m - nextPowTwo(m)
			startInd = i
			v = osamG.createVtx(id, Real, NONE, NONE, NONE, NONE)
			// fmt.Printf("Created vtx @ addr %v \n", v)
			inc_vtcs[i] = v
		} else { // Else: Edge encountered: continue in current tree run
			ii := i - startInd
			lvl := 0
			if ii > z {
				lvl = 1
			}
			// A. Create the actual leaf node
			me := osamG.createVtx(id, Inc, inpG.edges[i].U, NONE, NONE, NONE)
			// fmt.Printf("Created vtx @ addr %v \n", me)
			// fmt.Printf("My level: %v \n", lvl)
			inc_vtcs[i] = me
			if C[lvl] == NONE {
				C[lvl] = me
				// fmt.Printf("%v + here + lvl %v + %v \n", i, lvl, C[lvl])
			} else {
				parent := store[C[lvl]].UP
				store[me].UP = parent
				store[parent].RC = me
				C[lvl] = NONE
				// fmt.Printf("%v + here2 + %v \n", i, C[lvl])
			}
			// B. Create the corresponding internal node
			if ii == m {
				maxLvl := lg(m)
				store[C[maxLvl]].UP = v
				store[v].LC = C[maxLvl]
				C[maxLvl] = NONE
			} else {
				lvl = 1 + maxDivPowTwo(ii)
				if ii > z {
					lvl = 1 + maxDivPowTwo(2*ii-z)
				}
				intNode := osamG.createVtx(id, Internal, NONE, NONE, C[lvl-1], NONE)
				// fmt.Printf("Created vtx @ addr %v \n", intNode)
				// fmt.Printf("New level: %v \n", lvl)
				// fmt.Printf("LC node: %v \n", C[lvl-1])
				store[C[lvl-1]].UP = intNode
				if C[lvl] == NONE {
					C[lvl] = intNode
				} else {
					parent := store[C[lvl]].UP
					store[intNode].UP = parent
					store[parent].RC = intNode
					C[lvl] = NONE
				}
			}
		}
	}
	return inc_vtcs
}
